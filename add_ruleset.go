package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
)

var addRuleset = &cobra.Command{
	Use:   "add-ruleset",
	Short: "Adds a rector ruleset",
	Long:  `Edits the rector.php file to add a new ruleset`,
	Run: func(cmd *cobra.Command, args []string) {
		if !isRulesetArgumentValid(args) {
			log.Fatalf("invalid ruleset argument: %s. Example: \\Rector\\Set\\ValueObject\\LevelSetList::UP_TO_PHP_81\n", args[0])
		}
		file, lines, err := loadRectorFile()
		if err != nil {
			log.Fatalf(err.Error())
		}
		defer file.Close()

		rulesetInjectionPoint, err := findRulesetLineIndex(lines)
		if err != nil {
			log.Fatalf(err.Error())
		}

		lines = injectLine(lines, rulesetInjectionPoint, args[0])

		if err := writeRectorFile(file, lines); err != nil {
			log.Fatalf(err.Error())
		}
	},
}

func isRulesetArgumentValid(args []string) bool {
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

func findRulesetLineIndex(lines []string) (int, error) {
	for index, line := range lines {
		// are we on the line that represents a function call end?
		if !closingFunctionCallRegex.MatchString(line) {
			continue
		}

		// are we inside a ruleset function call, or another function call?
		for i := index; i >= 0; i-- {
			if strings.Contains(lines[i], "$rectorConfig->sets") {
				return index, nil
			}
		}
	}

	return 0, fmt.Errorf("failed finding ruleset section in rector.php")
}
