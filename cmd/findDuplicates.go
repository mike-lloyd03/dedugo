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
	"bytes"
	"fmt"
	"image"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	_ "github.com/adrium/goheif"
	"github.com/spf13/cobra"
	"github.com/vitali-fedulov/images"
)

// findDuplicatesCmd represents the findDuplicates command
var findDuplicatesCmd = &cobra.Command{
	Aliases: []string{"find", "f"},
	Args:    cobra.MinimumNArgs(2),
	Use:     "find-duplicates ref_directory eval_directory",
	Short:   "Finds duplicate images between two directories.",
	Long:    `Recursively searches through both input directories for images and compares if the "evaulation directory" contains any duplicates of images found in the "reference directory."`,
	Run: func(cmd *cobra.Command, args []string) {
		findDuplicates(args[0], args[1])
	},
}

func init() {
	rootCmd.AddCommand(findDuplicatesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// findDuplicatesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// findDuplicatesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

var (
	imgFormats = map[string]struct{}{
		".jpg":  {},
		".jpeg": {},
		".heic": {},
		".png":  {},
	}
	wg              sync.WaitGroup
	duplicatesFound int = 0
	m               sync.Mutex
)

type Image struct {
	Path string
	Hash []float32
	Size image.Point
}

type Pair struct {
	RefImage  string
	DupeImage string
}

func findDuplicates(refDir, evalDir string) {
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
	fmt.Printf("Done. %d potential duplicate images found.\n", duplicatesFound)
	// checkDuplicates(pairMap)
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

func getImagesFromDir(dir string) ([]Image, error) {
	imageList := make([]Image, 0)
	ch := make(chan Image, countFiles(dir))
	fmt.Println("Walking directory", dir)
	fmt.Println()

	filepath.WalkDir(dir, func(path string, entry fs.DirEntry, err error) error {
		fmt.Print("\033[1A\033[K")
		fmt.Println("Reading:", path)
		if !entry.IsDir() {
			if isImage(path) {
				wg.Add(1)
				go openAndHash(path, imageList, ch)
			}
		}
		return nil
	},
	)
	wg.Wait()
	close(ch)
	for img := range ch {
		imageList = append(imageList, img)
	}
	return imageList, nil
}

func openAndHash(path string, imageList []Image, ch chan Image) {
	defer wg.Done()
	img, err := OpenImage(path)
	if err != nil {
		log.Printf("Error opening %s: %s", path, err)
		return
	}
	hash, size := images.Hash(img)
	ch <- Image{path, hash, size}
}

// OpenImage opens and decodes an image file for a given path.
func OpenImage(path string) (img image.Image, err error) {
	fileBytes, err := os.ReadFile(path)
	if err != nil {
		return OpenImage(path)
	}
	file := bytes.NewReader(fileBytes)
	img, _, err = image.Decode(file)
	if err != nil {
		return nil, err
	}
	return img, err
}

func CompareImages(refImg Image, evalImages []Image, pairMap map[string]Pair) {
	defer wg.Done()
	for _, evalImg := range evalImages {
		if images.Similar(refImg.Hash, evalImg.Hash, refImg.Size, evalImg.Size) {
			m.Lock()
			duplicatesFound++
			pairMap[refImg.Path+","+evalImg.Path] = Pair{refImg.Path, evalImg.Path}
			m.Unlock()
		}
	}
}
