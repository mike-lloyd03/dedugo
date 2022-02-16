package cmd

import (
	"io"
	"io/ioutil"
	"log"
	"os"

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

func setupLogging(logToFile bool) {
	if logToFile {
		file, err := os.OpenFile("dedugo.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		// defer file.Close()
		if err != nil {
			log.Fatal(err)
		}
		log.SetOutput(file)
	} else {
		log.SetOutput(io.Discard)
	}
}
