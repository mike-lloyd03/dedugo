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
	"strings"
	"sync"
	"time"

	_ "github.com/adrium/goheif"
	"github.com/spf13/cobra"
	"github.com/vitali-fedulov/images"
)

var (
	results_path string
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
}

var (
	imgFormats = map[string]struct{}{
		".jpg":  {},
		".jpeg": {},
		".heic": {},
		".png":  {},
	}
	wg        sync.WaitGroup
	m         sync.Mutex
	openFiles int = 0
)

type Image struct {
	Path string
	Hash []float32
	Size image.Point
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
	// Setup logging
	file, err := os.OpenFile("dedugo.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)

	pairMap := make(map[string]Pair)

	refImages, err := getImagesFromDir(refDir)
	if err != nil {
		log.Fatal("Error:", err)
	}
	fmt.Printf("Images found in reference directory: %d images.\n", len(refImages))

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
}

func countFiles(dir string) int {
	count := 0
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		path := filepath.Join(dir, file.Name())
		if file.IsDir() {
			count += countFiles(path)
		} else if isImage(path) {
			count++
		}
	}
	return count
}

func isImage(path string) bool {
	_, found := imgFormats[strings.ToLower(filepath.Ext(path))]
	return found
}

func getImagesFromDir(dir string) (map[string]Image, error) {
	imageList := make(map[string]Image, 0)
	fmt.Println("Walking directory", dir)
	fmt.Println()

	filepath.WalkDir(dir, func(path string, entry fs.DirEntry, err error) error {
		fmt.Print("\033[1A\033[K")
		fmt.Println("Reading:", path)
		if !entry.IsDir() {
			if isImage(path) {
				wg.Add(1)
				go openAndHash(path, imageList)
			}
		}
		if err != nil {
			log.Fatal(path, err)
		}
		return nil
	},
	)
	wg.Wait()
	return imageList, nil
}

func openAndHash(path string, imageList map[string]Image) {
	defer wg.Done()
	img, err := OpenImage(path)
	if err != nil {
		// openAndHash(path, imageList)
		log.Fatalf("Error opening %s: %s", path, err)
	}
	hash, size := images.Hash(img)
	m.Lock()
	imageList[path] = Image{path, hash, size}
	m.Unlock()
}

// OpenImage opens and decodes an image file for a given path.
func OpenImage(path string) (img image.Image, err error) {
	if openFiles > 500 {
		time.Sleep(10 * time.Millisecond)
		return OpenImage(path)
	}
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	} else {
		m.Lock()
		openFiles++
		m.Unlock()
	}
	img, _, err = image.Decode(file)
	file.Close()
	m.Lock()
	openFiles--
	m.Unlock()
	if err != nil {
		return nil, err
	}
	return img, nil
}

func CompareImages(refImg Image, evalImages map[string]Image, pairMap map[string]Pair) {
	defer wg.Done()
	for _, evalImg := range evalImages {
		if images.Similar(refImg.Hash, evalImg.Hash, refImg.Size, evalImg.Size) {
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
