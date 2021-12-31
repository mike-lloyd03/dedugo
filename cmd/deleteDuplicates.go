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
	"os"

	"github.com/spf13/cobra"
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
}

func deleteDuplicates() {
	results := readResultsFile(results_path)
	for _, p := range results.ImagePairs {
		if p.Confirmed {
			fmt.Println("Deleting", p.DupeImage)
			os.Remove(p.DupeImage)
		}
	}
	fmt.Println("Done.")
}
