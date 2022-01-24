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
	"io"
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
	results_path string
	logToFile    bool
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

	findDuplicatesCmd.Flags().StringVarP(&results_path, "output", "o", "dedugo_results.yaml", "output file for results")
	findDuplicatesCmd.Flags().BoolVar(&logToFile, "log", false, "log events to file")
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
	RefImage  string `yaml:"ReferenceImage"`
	DupeImage string `yaml:"DuplicateImage"`
	Confirmed bool   `yaml:"Confirmed?"`
}

type Results struct {
	RefDir     string `yaml:"ReferenceDirectory"`
	EvalDir    string `yaml:"EvaluationDirectory"`
	StartIdx   int    `yaml:"StartIndex"`
	ImagePairs []Pair `yaml:"ImagePairs"`
}

func findDuplicates(refDir, evalDir string) {
	startTime := time.Now()

	// Setup logging
	if logToFile {
		file, err := os.OpenFile("dedugo.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		defer file.Close()
		if err != nil {
			log.Fatal(err)
		}
		log.SetOutput(file)
	} else {
		log.SetOutput(io.Discard)
	}

	log.Printf("Finding duplicates for %s and %s\n", refDir, evalDir)

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
	log.Printf("Done. Total elapsed time: %s", time.Now().Sub(startTime).Round(10*time.Millisecond))
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

	var numWorkers int
	if len(imgs) < maxWorkers {
		numWorkers = len(imgs)
	} else {
		numWorkers = maxWorkers
	}
	startTime := time.Now()
	log.Printf("Beginning scan of %s with %d workers.\n", dir, numWorkers)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go openAndHashWorker(pathChan, imageChan)
	}

	for _, path := range imgs {
		pathChan <- path
	}

	close(pathChan)
	wg.Wait()
	close(imageChan)
	for img := range imageChan {
		imageList = append(imageList, img)
	}
	log.Printf("Finished scan. Found %d images. Elapsed time: %s\n", len(imageList), time.Now().Sub(startTime).Round(10*time.Millisecond))
	return imageList, nil
}

func openAndHashWorker(pathChan <-chan string, imageChan chan<- Image) {
	defer wg.Done()
	for path := range pathChan {
		img, err := OpenImage(path)
		if err != nil {
			log.Fatalf("Error opening %s: %s", path, err)
		}
		icon := images.Icon(img, path)
		imageChan <- Image{path, icon}
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
		if images.Similar(refImg.Icon, evalImg.Icon) {
			m.Lock()
			pairMap[refImg.Path+","+evalImg.Path] = Pair{RefImage: refImg.Path, DupeImage: evalImg.Path}
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
	WriteResultsFile(results, results_path)
}
