package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var addRuleset = &cobra.Command{
	Use:   "add-ruleset",
	Short: "Adds a rector ruleset",
	Long:  `Edits the rector.php file to add a new ruleset`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: implement this validation:
		//if !isRulesetArgumentValid(args) {
		//	log.Fatalf("invalid ruleset argument: %s. Example: \\Rector\\Set\\ValueObject\\LevelSetList::UP_TO_PHP_81\n", args[0])
		//}
		const rectorFile = "rector.php"
		file, err := os.OpenFile(rectorFile, os.O_RDWR, 0644)
		if err != nil {
			log.Fatalf("failed opening rector.php file: %s", err)
		}
		defer file.Close()

		lines, err := linesFromReader(file)
		if err != nil {
			log.Fatalf("failed reading rector.php file: %s", err)
		}

		for index, line := range lines {
			if rulesetInjectedSuccessfully(line, index, lines, args) {
				break
			}
			if index == len(lines)-1 {
				log.Fatalf("failed injecting ruleset: %s", args[0])
			}
		}

		for _, line := range lines {
			log.Println(line)
		}
	},
}

func rulesetInjectedSuccessfully(line string, index int, lines []string, args []string) bool {
	if !lineAtTheEndOfRuleset(line, index, lines) {
		return false
	}

	lastRulesetLine := lines[index-1]

	indentSize := strings.Count(line, " ") + 4

	if !strings.HasSuffix(lastRulesetLine, ",") {
		if strings.Contains(lastRulesetLine, "::") {
			lines[index-1] = lastRulesetLine + ","
		}
	}
	lines = append(
		lines[:index],
		append(
			[]string{fmt.Sprintf("%s%s,", strings.Repeat(" ", indentSize), args[0])},
			lines[index:]...,
		)...,
	)

	return true
}

func lineAtTheEndOfRuleset(line string, index int, lines []string) bool {
	closingFunctionCallRegex, err := regexp.Compile(`]?\)?;`)
	if err != nil {
		log.Fatalf("failed compiling regex: %s", err)
	}

	if !closingFunctionCallRegex.MatchString(line) {
		return false
	}

	for i := index; i >= 0; i-- {
		if strings.Contains(lines[i], "$rectorConfig->sets") {
			break
		}
		if i == 0 {
			return false
		}
	}

	return true
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
