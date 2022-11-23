package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var bumpPhp = &cobra.Command{
	Use:   "bump-php",
	Short: "Bump PHP version",
	Long:  `Bump PHP version in composer.json. Must be ran in the root of the project, where the composer.json file is located`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: fetch latest/requested php version from the php repo website thingy
		phpVersion := "8.2"
		if len(args) > 0 {
			phpVersion = args[0]
		}

		// Try to open the composer.json file
		const pathToComposerJson = "./composer.json"
		file, err := os.OpenFile(pathToComposerJson, os.O_RDWR, 0644)
		if err != nil {
			log.Fatalf("missing file: %s %v\n", pathToComposerJson, err)
		}

		// Decode the json file into a struct
		var s schema
		if err := json.NewDecoder(file).Decode(&s); err != nil {
			log.Fatalf("failed decoding composer.json: %v\n", err)
		}

		s.setPhpVersion(phpVersion)

		// Write the updated struct back to the file
		if err := file.Truncate(0); err != nil {
			log.Fatalf("failed truncating %s %v\n", pathToComposerJson, err)
		}

		if _, err := file.Seek(0, 0); err != nil {
			return
		}
		if err := json.NewEncoder(file).Encode(s); err != nil {
			log.Fatalf("failed encoding composer.json: %v\n", err)
		}
	},
}
