package main

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
	"github.com/vitali-fedulov/images/v2"
)

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

func main() {
	pairMap := make(map[string]Pair)
	refDir := os.Args[1]
	evalDir := os.Args[2]

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
	checkDuplicates(pairMap)
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
