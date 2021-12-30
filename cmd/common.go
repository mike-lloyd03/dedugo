package cmd

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

func readResultsFile(path string) Results {
	results := Results{}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("Error reading results file.", err)
	}
	yaml.Unmarshal(file, &results)
	return results
}

func WriteResultsFile(results Results, path string) {
	data, err := yaml.Marshal(results)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(path, data, 0644)
	if err != nil {
		log.Fatal("Error writing results file.", err)
	}
}
