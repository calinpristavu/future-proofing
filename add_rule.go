package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
)

var addRule = &cobra.Command{
	Use:   "add-rule",
	Short: "Adds a rector rule",
	Long:  `Edits the rector.php file to add a new rule`,
	Run: func(cmd *cobra.Command, args []string) {
		if !isRuleArgumentValid(args) {
			log.Fatalf("invalid ruleset argument: %s. Example: \\Rector\\Set\\ValueObject\\LevelSetList::UP_TO_PHP_81\n", args[0])
		}

		file, lines, err := loadRectorFile()
		if err != nil {
			log.Fatalf(err.Error())
		}
		defer file.Close()

		ruleInjectionPoint, err := findRuleLineIndex(lines)
		if err != nil {
			// if we can't find a ->rules([...]) section, we'll try to find a ->rule(...) section and convert it to a ->rules section
			lines = convertSingleRuleToMultipleRules(lines)

			ruleInjectionPoint, err = findRuleLineIndex(lines)
			if err != nil {
				log.Fatalf(err.Error())
			}
		}

		lines = injectLine(lines, ruleInjectionPoint, args[0])

		if err := writeRectorFile(file, lines); err != nil {
			log.Fatalf(err.Error())
		}
	},
}

func convertSingleRuleToMultipleRules(lines []string) []string {
	const singleRuleString = "$rectorConfig->rule("
	for index, line := range lines {
		if !strings.Contains(line, "$rectorConfig->rule(") {
			continue
		}

		indentation := line[:strings.Index(line, singleRuleString)]

		lines[index] = fmt.Sprintf("%s$rectorConfig->rules([", indentation)
		lines = append(
			lines[:index+1],
			append(
				[]string{fmt.Sprintf("%s)];", indentation)},
				lines[index+1:]...,
			)...,
		)

		if strings.Contains(lines[index-1], "register a single rule") {
			lines[index-1] = fmt.Sprintf("%s// register multiple rules", indentation)
		}
	}

	return lines
}

func isRuleArgumentValid(args []string) bool {
	// there should be an argument
	if len(args) == 0 {
		return false
	}

	// argument must be a call to a constant
	if !strings.Contains(args[0], "::") {
		return false
	}

	// argument must have a namespace
	if !strings.Contains(args[0], "\\") {
		return false
	}

	return true
}

func findRuleLineIndex(lines []string) (int, error) {
	for index, line := range lines {
		// are we on the line that represents a function call end?
		if !closingFunctionCallRegex.MatchString(line) {
			continue
		}

		// are we inside a ruleset function call, or another function call?
		for i := index; i >= 0; i-- {
			if strings.Contains(lines[i], "$rectorConfig->rules") {
				return index, nil
			}
		}
	}

	return 0, fmt.Errorf("failed finding ruleset section in rector.php")
}
