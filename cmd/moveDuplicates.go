/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var moveAll bool

// moveDuplicatesCmd represents the moveDuplicates command
var moveDuplicatesCmd = &cobra.Command{
	Aliases: []string{"move", "m"},
	Args:    cobra.MinimumNArgs(1),
	Use:     "move-duplicates destination_directory",
	Short:   "Move all confirmed duplicates to a designated directory",
	Long:    `After running "dedugo find-duplicates" and "dedugo check-results", this command will move all confirmed duplicate images to the designated directory.`,
	Run: func(cmd *cobra.Command, args []string) {
		moveDuplicates(args[0])
	},
}

func init() {
	rootCmd.AddCommand(moveDuplicatesCmd)

	moveDuplicatesCmd.Flags().StringVarP(&resultsPath, "input-file", "i", "dedugo_results.yaml", "input file to read results from")
	moveDuplicatesCmd.Flags().BoolVar(&moveAll, "all", false, "move all duplicate images whether they are confirmed or not")
	moveDuplicatesCmd.Flags().BoolVar(&logToFile, "log", false, "log events to file")
	moveDuplicatesCmd.Flags().BoolVar(&dryRun, "dry-run", false, "show only what would be deleted without actually doing it")
}

func moveDuplicates(destDir string) {
	setupLogging(logToFile)

	var input string
	fmt.Printf("Are you sure you want to move all duplicate images found in %s? [y/N]: ", resultsPath)
	fmt.Scan(&input)
	if strings.ToLower(input) != "yes" && strings.ToLower(input) != "y" {
		fmt.Println("Aborting")
		return
	}

	results := readResultsFile(resultsPath)
	log.Printf("Moving duplicate images to %s.", destDir)
	for _, p := range results.ImagePairs {
		if p.Confirmed || moveAll {
			newPath := filepath.Join(destDir, filepath.Base(p.DupeImage))
			log.Printf("Moving %s to %s\n", p.DupeImage, newPath)
			if !dryRun {
				err := os.Rename(p.DupeImage, newPath)
				if err != nil {
					log.Printf("%s", err)
				}
			}
		}
	}
	fmt.Println("Done.")
}
