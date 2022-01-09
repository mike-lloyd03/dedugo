package imagelist

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"testing"
)

func TestNewImageList(t *testing.T) {
	paths := []string{
		"./test_images/Obi1.jpg",
		"./test_images/Obi2.jpg",
		"./test_images/Jango1.jpg",
	}

	imgList := make([]image.Image, 3)
	for i, path := range paths {
		imgFile, err := os.Open(path)
		if err != nil {
			t.Errorf("failed to open %s %s", path, err)
		}
		img, _, err := image.Decode(imgFile)
		if err != nil {
			t.Errorf("failed to decode %s %s", path, err)
		}
		imgFile.Close()
		imgList[i] = img
	}

	expectIL := ImageList{
		index:      0,
		paths:      paths,
		imageCache: imgList,
	}

	gotIL, err := New(paths)
	if err != nil {
		t.Error("could not create new imageList type")
	}

	// Compare index
	if gotIL.index != 0 {
		t.Error("index field was not instantiated correctly")
	}

	// Compare paths
	for i, p := range gotIL.paths {
		if p != expectIL.paths[i] {
			t.Error("path field was not instantiated correctly")
		}
	}

	// Compare cached images
	for i, img := range gotIL.imageCache {
		if !compareImages(img, gotIL.imageCache[i]) {
			t.Error("image cache was not instantiated correctly")
		}
	}
}

func TestNextandPrevious(t *testing.T) {
	// Setup
	paths := []string{
		"./test_images/Obi1.jpg",
		"./test_images/Obi2.jpg",
		"./test_images/Jango1.jpg",
		"./test_images/Jango2.jpg",
		"./test_images/Kylo1.jpg",
	}
	expectImages := make([]image.Image, len(paths))
	for i, p := range paths {
		img, _ := loadImage(p)
		expectImages[i] = img
	}

	// ----------
	// Test Next
	// ----------
	il, _ := New(paths)
	nextImage, err := il.Next()
	if err != nil {
		t.Errorf("could not get next (first) image. %s", err)
	}

	// nextImage should be the first image in the list
	if !compareImages(nextImage, expectImages[0]) {
		t.Error("got wrong next (first) image")
	}

	// Index should be 1
	if il.index != 1 {
		t.Error("index was not iterated. index =", il.index)
	}

	// Image cache should contain the first three images in the list
	if len(il.imageCache) != 3 {
		t.Error("image cache size should be 3. is ", len(il.imageCache))
	}
	for i, img := range il.imageCache {
		if !compareImages(img, expectImages[i]) {
			t.Error("image cache is not correct at index", i)
		}
	}
	// Get second image
	nextImage, err = il.Next()
	if err != nil {
		t.Errorf("could not get next (second) image. %s", err)
	}
	// nextImage should be the first image in the list
	if !compareImages(nextImage, expectImages[1]) {
		t.Error("got wrong next (second) image")
	}

	// Get third image
	nextImage, err = il.Next()
	if err != nil {
		t.Errorf("could not get next (third) image. %s", err)
	}
	// nextImage should be the third image in the list
	if !compareImages(nextImage, expectImages[2]) {
		t.Error("got wrong next (third) image")
	}

	// First image in the image cache should be the second image in the list
	if !compareImages(il.imageCache[0], expectImages[1]) {
		t.Error("first image in image cache should be second image in path list")
	}

	// Last image in the image cache should be the fourth image in the list
	if !compareImages(il.imageCache[len(il.imageCache)-1], expectImages[3]) {
		t.Error("last image in the image cache should be fourth image in the path list")
	}

	// Get fourth image
	nextImage, err = il.Next()
	if err != nil {
		t.Errorf("could not get next (fourth) image. %s", err)
	}
	// nextImage should be the fourth image in the list
	if !compareImages(nextImage, expectImages[3]) {
		t.Error("got wrong next (fourth) image")
	}
	// First image in image cache should be second to last image in list
	if !compareImages(il.imageCache[0], expectImages[2]) {
		t.Error("first image in image cache should be second to last image in path list")
	}
	// Last image in image cache should be last image in list
	if !compareImages(il.imageCache[2], expectImages[4]) {
		t.Error("last image in image cache should be last image in path list")
	}

	// Get fifth image
	nextImage, err = il.Next()
	if err != nil {
		t.Errorf("could not get next (fifth) image. %s", err)
	}
	// nextImage should be the fifth image in the list
	if !compareImages(nextImage, expectImages[4]) {
		t.Error("got wrong next (fifth) image")
	}

	// Getting next image beyond length of list should return the last image
	nextImage, err = il.Next()
	if err != nil {
		t.Errorf("could not get next (last) image after index is max. %s", err)
	}
	if !compareImages(nextImage, expectImages[4]) {
		t.Error("got wrong next (last) image after index is max")
	}
	// Index should not be greater than len(paths) - 1
	if il.index != len(il.paths)-1 {
		t.Error("index should be len(paths) - 1")
	}

	// ----------
	// Test Previous
	// ----------
	// Get fourth image
	nextImage, err = il.Previous()
	if err != nil {
		t.Errorf("could not get previous (fourth) image. %s", err)
	}
	// nextImage should be the fourth image in the list
	if !compareImages(nextImage, expectImages[3]) {
		t.Error("got wrong previous (fourth) image")
	}

	// First image in image cache should be second to last image in list
	if !compareImages(il.imageCache[0], expectImages[2]) {
		t.Error("first image in image cache should be second to last image in path list")
	}
	// Last image in image cache should be last image in list
	if !compareImages(il.imageCache[2], expectImages[4]) {
		t.Error("last image in image cache should be last image in path list")
	}

	// Get third image
	prevImage, err := il.Previous()
	if err != nil {
		t.Errorf("could not get previous (third) image. %s", err)
	}
	// prevImage should be the fourth image in the list
	if !compareImages(prevImage, expectImages[2]) {
		t.Error("got wrong previous (third) image")
	}
	// Get second image
	prevImage, err = il.Previous()
	if err != nil {
		t.Errorf("could not get previous (second) image. %s", err)
	}
	// prevImage should be the third image in the list
	if !compareImages(prevImage, expectImages[1]) {
		t.Error("got wrong previous (second) image")
	}
	// Get first image
	prevImage, err = il.Previous()
	if err != nil {
		t.Errorf("could not get previous (first) image. %s", err)
	}
	// prevImage should be the first image in the list
	if !compareImages(prevImage, expectImages[0]) {
		t.Error("got wrong previous (first) image. got", findImage(prevImage, expectImages))
	}
	// Getting previous image before zero should return the first image
	prevImage, err = il.Previous()
	if err != nil {
		t.Errorf("could not get previous (first) image after index is zero. %s", err)
	}
	// prevImage should be the first image in the list
	if !compareImages(prevImage, expectImages[0]) {
		t.Error("got wrong previous (first) image after index is zero")
	}

	// Index should not be less than zero
	if il.index != 0 {
		t.Error("index should be zero")
	}
}

func compareImages(img1, img2 image.Image) bool {
	if img1.Bounds() != img2.Bounds() {
		return false
	}
	r := img1.Bounds()
	for x := r.Min.X; x < r.Max.X; x++ {
		for y := r.Min.Y; y < r.Max.Y; y++ {
			if img1.At(x, y) != img2.At(x, y) {
				return false
			}
		}
	}
	return true
}

func findImage(img image.Image, imgList []image.Image) int {
	for i, compImage := range imgList {
		if compareImages(img, compImage) {
			return i
		}
	}
	return 100
}
