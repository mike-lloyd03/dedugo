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
	"os/exec"

	"github.com/spf13/cobra"
)

// checkResultsCmd represents the checkResults command
var checkResultsCmd = &cobra.Command{
	Aliases: []string{"check", "c"},
	Use:     "check-results",
	Short:   "Check each of the image pairs found in the \"find-duplicates\" command",
	Long:    `Check each of the image pairs by opening both of them in the system default image application. The user will be prompted to confirm if the file is a duplicate or not. All confirmed duplicates can subsequently be deleted with the "delete" command.`,
	Run: func(cmd *cobra.Command, args []string) {
		checkResultsGui()
	},
}

func init() {
	rootCmd.AddCommand(checkResultsCmd)

	checkResultsCmd.Flags().StringVarP(&results_path, "input", "i", "dedugo_results.yaml", "input file to read results from")
}

// checkResults reads from the Results file and iterates over the Image Pairs,
// asking the user to confirm if the image is a duplicate or not
func checkResults(results_path string) {
	var input string
	results := readResultsFile(results_path)

read_input:
	for i := results.StartIdx; i < len(results.ImagePairs); i++ {
		p := results.ImagePairs[i]
		openDuplicates(p.RefImage, p.DupeImage)

		results.StartIdx = i
		WriteResultsFile(results, results_path)

		fmt.Printf("%s and %s are duplicates? [y/N/stop] ", p.RefImage, p.DupeImage)
		fmt.Scanln(&input)

		switch input {
		case "y", "Y":
			results.ImagePairs[i].Confirmed = true
			WriteResultsFile(results, results_path)
		case "stop":
			break read_input
		default:
			continue
		}
		// if gui.ShowGui(p.RefImage, p.DupeImage) {
		// 	fmt.Println("yes")
		// 	results.ImagePairs[i].Confirmed = true
		// 	writeResultsFile(results, results_path)
		// }
	}
}

func openDuplicates(refFile string, dupeFile string) {
	exec.Command("xdg-open", refFile).Start()
	exec.Command("xdg-open", dupeFile).Run()
}

func checkResultsGui() {
	showGui()
}
