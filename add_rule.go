package main

import (
	"github.com/spf13/cobra"
)

var addRule = &cobra.Command{
	Use:   "add-rule",
	Short: "Adds a rector rule",
	Long:  `Edits the rector.php file to add a new rule`,
	Run: func(cmd *cobra.Command, _ []string) {

	},
}
