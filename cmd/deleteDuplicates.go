/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	"strings"

	"github.com/spf13/cobra"
)

var (
	deleteAll bool
	dryRun    bool
)

// deleteDuplicatesCmd represents the deleteDuplicates command
var deleteDuplicatesCmd = &cobra.Command{
	Aliases: []string{"delete", "d"},
	Use:     "delete-duplicates",
	Short:   "Delete all confirmed duplicate files",
	Long:    `After running "dedugo find-duplicates" and "dedugo check-results", this command will go through the file system and delete all confirmed duplicate images.`,
	Run: func(cmd *cobra.Command, args []string) {
		deleteDuplicates()
	},
}

func init() {
	rootCmd.AddCommand(deleteDuplicatesCmd)

	deleteDuplicatesCmd.Flags().StringVarP(&resultsPath, "input-file", "i", "dedugo_results.yaml", "input file to read results from")
	deleteDuplicatesCmd.Flags().BoolVar(&deleteAll, "all", false, "delete all duplicate images whether they are confirmed or not")
	deleteDuplicatesCmd.Flags().BoolVar(&logToFile, "log", false, "log events to file")
	deleteDuplicatesCmd.Flags().BoolVar(&dryRun, "dry-run", false, "show only what would be deleted without actually doing it")
}

func deleteDuplicates() {
	setupLogging(logToFile)

	var input string
	fmt.Printf("Are you sure you want to delete all duplicate images found in %s? [y/N]: ", resultsPath)
	fmt.Scan(&input)
	if strings.ToLower(input) != "yes" && strings.ToLower(input) != "y" {
		fmt.Println("Aborting")
		return
	}

	results := readResultsFile(resultsPath)
	log.Printf("Deleting duplicate images.")
	for _, p := range results.ImagePairs {
		if p.Confirmed || deleteAll {
			fmt.Println("Deleting", p.DupeImage)
			if !dryRun {
				err := os.Remove(p.DupeImage)
				if err != nil {
					log.Printf("Failed to delete %s. %s\n", p.DupeImage, err)
				}
			}
		}
	}
	fmt.Println("Done.")
}
