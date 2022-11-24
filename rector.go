package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

var (
	closingFunctionCallRegex *regexp.Regexp
	tabIndentationRegex      *regexp.Regexp
	spaceIndentationRegex    *regexp.Regexp
)

func init() {
	var err error
	closingFunctionCallRegex, err = regexp.Compile(`]?\)?;`)
	if err != nil {
		log.Fatalf("failed compiling closingFunctionCallRegex: %s", err)
	}

	tabIndentationRegex, err = regexp.Compile(`^\t+`)
	if err != nil {
		log.Fatalf("failed compiling tabIndentationRegex: %s", err)
	}

	spaceIndentationRegex, err = regexp.Compile(`^\s+`)
	if err != nil {
		log.Fatalf("failed compiling spaceIndentationRegex: %s", err)
	}
}

func injectLine(lines []string, index int, line string) []string {
	lines[index-1] = ensureLineHasComma(lines[index-1])

	line = ensureGlobalNamespace(line)

	indentation := computeIndentation(lines[index])
	lines = append(
		lines[:index],
		append(
			[]string{fmt.Sprintf("%s%s,", indentation, line)},
			lines[index:]...,
		)...,
	)

	return lines
}

func ensureGlobalNamespace(line string) string {
	if strings.HasPrefix(line, "\\") {
		return line
	}

	return "\\" + line
}

func ensureLineHasComma(line string) string {
	if strings.HasSuffix(line, ",") {
		return line
	}

	if strings.Contains(line, "::") {
		return line + ","
	}

	return line
}

func computeIndentation(line string) string {
	indentation := tabIndentationRegex.FindString(line)
	if indentation != "" {
		return strings.Repeat("\t", strings.Count(indentation, "\t")+1)
	}

	indentation = spaceIndentationRegex.FindString(line)
	if indentation != "" {
		return strings.Repeat(" ", strings.Count(indentation, " ")+4)
	}

	return ""
}

func writeRectorFile(file *os.File, lines []string) error {
	if _, err := file.Seek(0, 0); err != nil {
		return fmt.Errorf("failed seeking to the beginning of the file: %s", err)
	}

	if err := file.Truncate(0); err != nil {
		return fmt.Errorf("failed truncating the file: %s", err)
	}

	for _, line := range lines {
		if _, err := file.WriteString(line + "\n"); err != nil {
			return fmt.Errorf("failed writing line to file: %s", err)
		}
	}

	return nil
}

func loadRectorFile() (*os.File, []string, error) {
	const rectorFile = "rector.php"

	file, err := os.OpenFile(rectorFile, os.O_RDWR, 0644)
	if err != nil {
		return nil, nil, fmt.Errorf("failed opening rector.php file: %s", err)
	}

	lines, err := linesFromReader(file)
	if err != nil {
		return nil, nil, fmt.Errorf("failed reading rector.php file: %s", err)
	}

	return file, lines, nil
}

func linesFromReader(r io.Reader) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
