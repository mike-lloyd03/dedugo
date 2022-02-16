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
	"image"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	_ "github.com/adrium/goheif"
	"github.com/spf13/cobra"
	images "github.com/vitali-fedulov/images3"
)

var (
	resultsPath   string
	logToFile     bool
	minConfidence int
)

// findDuplicatesCmd represents the findDuplicates command
var findDuplicatesCmd = &cobra.Command{
	Aliases: []string{"find", "f"},
	Args:    cobra.MinimumNArgs(2),
	Use:     "find-duplicates ref_directory eval_directory",
	Short:   "Finds duplicate images between two directories.",
	Long:    `Recursively searches through both input directories for images and compares if the "evaulation directory" contains any duplicates of images found in the "reference directory".`,
	Run: func(cmd *cobra.Command, args []string) {
		findDuplicates(args[0], args[1])
	},
}

func init() {
	rootCmd.AddCommand(findDuplicatesCmd)

	findDuplicatesCmd.Flags().StringVarP(&resultsPath, "output-file", "o", "dedugo_results.yaml", "output file for results")
	findDuplicatesCmd.Flags().BoolVar(&logToFile, "log", false, "log events to file")
	findDuplicatesCmd.Flags().IntVarP(&minConfidence, "min-confidence", "m", 1, "set the minimum confidence score (1-5) required to consider images similar")

	if minConfidence < 1 || minConfidence > 5 {
		log.Fatal("Minimum confidence must be in the range of 1-5")
	}
}

var (
	imgFormats = map[string]struct{}{
		".jpg":  {},
		".jpeg": {},
		".heic": {},
		".png":  {},
	}
	wg         sync.WaitGroup
	m          sync.Mutex
	maxWorkers int
)

type Image struct {
	Path string
	Icon images.IconT
}

type Pair struct {
	RefImage   string `yaml:"ReferenceImage"`
	DupeImage  string `yaml:"DuplicateImage"`
	Confirmed  bool   `yaml:"Confirmed?"`
	Confidence int    `yaml:"Confidence"`
}

type Results struct {
	RefDir     string `yaml:"ReferenceDirectory"`
	EvalDir    string `yaml:"EvaluationDirectory"`
	StartIdx   int    `yaml:"StartIndex"`
	ImagePairs []Pair `yaml:"ImagePairs"`
}

func findDuplicates(refDir, evalDir string) {
	startTime := time.Now()

	setupLogging(logToFile)

	log.Printf("Finding duplicates for %s and %s. Minimum confidence score = %d.\n", refDir, evalDir, minConfidence)

	maxWorkers = runtime.NumCPU()

	pairMap := make(map[string]Pair)

	fmt.Println("Walking reference directory", refDir)
	refImages, err := getImagesFromDir(refDir)
	if err != nil {
		log.Fatal("Error:", err)
	}
	fmt.Printf("Images found in reference directory: %d images.\n", len(refImages))

	fmt.Println("Walking evaluation directory", evalDir)
	evalImages, err := getImagesFromDir(evalDir)
	if err != nil {
		log.Fatal("Error:", err)
	}
	fmt.Printf("Images found in evaluation directory: %d images.\n", len(evalImages))

	fmt.Println("Comparing images...")
	for _, refImg := range refImages {
		wg.Add(1)
		go CompareImages(refImg, evalImages, pairMap)
	}
	wg.Wait()
	fmt.Printf("Done. %d potential duplicate images found.\n", len(pairMap))
	// checkDuplicates(pairMap)
	GenerateResults(refDir, evalDir, pairMap)
	log.Printf("Done. Found %d potential duplicates. Total elapsed time: %s", len(pairMap), time.Now().Sub(startTime).Round(10*time.Millisecond))
}

func getImagePaths(dir string) []string {
	paths := make([]string, 0)
	filepath.WalkDir(dir, func(path string, entry fs.DirEntry, err error) error {
		if !entry.IsDir() {
			if isImage(path) {
				paths = append(paths, path)
			}
		}
		if err != nil {
			log.Fatal(path, err)
		}
		return nil
	},
	)
	return paths
}

func isImage(path string) bool {
	_, found := imgFormats[strings.ToLower(filepath.Ext(path))]
	return found
}

func getImagesFromDir(dir string) ([]Image, error) {
	imgs := getImagePaths(dir)
	imageList := make([]Image, 0)
	pathChan := make(chan string, len(imgs))
	imageChan := make(chan Image, len(imgs))
	countChan := make(chan bool, len(imgs))

	var numWorkers int
	if len(imgs) < maxWorkers {
		numWorkers = len(imgs)
	} else {
		numWorkers = maxWorkers
	}
	startTime := time.Now()
	log.Printf("Beginning scan of %s with %d workers.\n", dir, numWorkers)

	go monitorProgress(countChan, len(imgs))

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go openAndHashWorker(pathChan, imageChan, countChan)
	}

	for _, path := range imgs {
		pathChan <- path
	}

	close(pathChan)
	wg.Wait()
	close(imageChan)
	close(countChan)
	for img := range imageChan {
		imageList = append(imageList, img)
	}
	log.Printf("Finished scan. Found %d images. Elapsed time: %s\n", len(imageList), time.Now().Sub(startTime).Round(10*time.Millisecond))
	return imageList, nil
}

func openAndHashWorker(pathChan <-chan string, imageChan chan<- Image, countChan chan<- bool) {
	defer wg.Done()
	for path := range pathChan {
		img, err := OpenImage(path)
		if err != nil {
			log.Fatalf("Error opening %s: %s", path, err)
		}
		icon := images.Icon(img, path)
		imageChan <- Image{path, icon}
		countChan <- true
	}
}

// OpenImage opens and decodes an image file for a given path.
func OpenImage(path string) (img image.Image, err error) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	img, _, err = image.Decode(file)
	file.Close()
	if err != nil {
		return nil, err
	}
	return img, nil
}

func CompareImages(refImg Image, evalImages []Image, pairMap map[string]Pair) {
	defer wg.Done()
	for _, evalImg := range evalImages {
		// m1, m2, m3 := images.EucMetric(refImg.Icon, evalImg.Icon)
		// if (m1+m2+m3)/3 < 2000 {
		// 	m.Lock()
		// 	pairMap[refImg.Path+","+evalImg.Path] = Pair{RefImage: refImg.Path, DupeImage: evalImg.Path}
		// 	m.Unlock()
		// }
		confidence := calcConfidence(images.EucMetric(refImg.Icon, evalImg.Icon))
		if confidence >= minConfidence {
			m.Lock()
			pairMap[refImg.Path+","+evalImg.Path] = Pair{RefImage: refImg.Path, DupeImage: evalImg.Path, Confidence: confidence}
			m.Unlock()
		}
	}
}

func GenerateResults(refDir, evalDir string, pairMap map[string]Pair) {
	pairArray := make([]Pair, len(pairMap))
	i := 0
	for _, p := range pairMap {
		pairArray[i] = p
		i++
	}
	results := Results{
		RefDir:     refDir,
		EvalDir:    evalDir,
		StartIdx:   0,
		ImagePairs: pairArray,
	}
	WriteResultsFile(results, resultsPath)
}

// calcConfidence returns the confidence that two images are similar on a scale of 0-5
// with 5 being the highest confidence
func calcConfidence(m1, m2, m3 float32) int {
	avg := (m1 + m2 + m3) / 3
	if avg < 2000 {
		return 5
	}
	if avg < 5000 {
		return 4
	}
	if avg < 8000 {
		return 3
	}
	if avg < 11000 {
		return 2
	}
	if avg < 14000 {
		return 1
	}
	return 0
}

func monitorProgress(countChan chan bool, total int) {
	count := 1
	fmt.Println()
	for done := range countChan {
		if done {
			fmt.Print("\033[1A\033[K")
			fmt.Printf("Scanning in progress: %d / %d\n", count, total)
			count += 1
		}
	}
}
