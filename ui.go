package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
)

func checkDuplicates(pairMap map[string]Pair) {
	var input string
	outputFile, err := os.Create("./duplicateImages.txt")
	if err != nil {
		log.Fatal(err)
	}
	writer := bufio.NewWriter(outputFile)
	defer writer.Flush()

input:
	for _, p := range pairMap {
		openDuplicates(p.RefImage, p.DupeImage)

		fmt.Printf("%s and %s are duplicates? [y/N] ", p.RefImage, p.DupeImage)
		fmt.Scanln(&input)
		switch input {
		case "y", "Y":
			_, err := writer.WriteString(p.DupeImage + "\n")
			if err != nil {
				log.Fatal(err)
			}
		case "stop":
			break input
		default:
			continue
		}
	}
}

func openDuplicates(refFile string, dupeFile string) {
	exec.Command("xdg-open", refFile).Start()
	exec.Command("xdg-open", dupeFile).Run()
}
