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
		"./test_images/Jango3.jpg",
		"./test_images/Jango4.jpg",
	}

	imgList := make([]image.Image, 3)
	for i, path := range paths[:3] {
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
		if !compareImages(img, imgList[i]) {
			t.Error("image cache was not instantiated correctly")
		}
	}

	// Load image list with less than 3 items
	_, err = New(paths[:2])
	if err == nil {
		t.Error("creating an ImageList with less than 3 items should return an error")
	}
	// Load image list with un-loadable image
	badPaths := append(paths[:2], "notANimagePath.jpg")
	_, err = New(badPaths)
	if err == nil {
		t.Error("creating an ImageList with unopenable images should return an error")
	}

	// Load image list at index > 0
	gotIL, err = New(paths, 1)
	if err != nil {
		t.Log(err)
		t.Error("failed to create image list at index > 0")
	}

	// Compare index
	if gotIL.index != 1 {
		t.Error("index field for new image list index > 0 was not instantiated correctly")
	}

	// Compare paths
	for i, p := range gotIL.paths {
		if p != expectIL.paths[i] {
			t.Error("path field was not instantiated correctly")
		}
	}

	// Compare cached images
	for i, img := range gotIL.imageCache {
		if !compareImages(img, imgList[i+1]) {
			t.Error("image cache was not instantiated correctly")
		}
	}

}

func TestNextandPrevious(t *testing.T) {
	// Setup
	paths := []string{
		"./test_images/Obi1.jpg",
		"./test_images/Obi2.jpg",
		"./test_images/Jango3.jpg",
		"./test_images/Jango4.jpg",
		"./test_images/Kylo5.jpg",
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
	nextImage, nextPath, err := il.Next()
	if err != nil {
		t.Errorf("could not get next (first) image. %s", err)
	}
	if nextPath != paths[0] {
		t.Error("got wrong next (first) image path")
	}
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
	nextImage, nextPath, err = il.Next()
	if err != nil {
		t.Errorf("could not get next (second) image. %s", err)
	}
	if nextPath != paths[1] {
		t.Error("got wrong next (second) image path")
	}
	if !compareImages(nextImage, expectImages[1]) {
		t.Error("got wrong next (second) image")
	}

	// Get third image
	nextImage, nextPath, err = il.Next()
	if err != nil {
		t.Errorf("could not get next (third) image. %s", err)
	}
	if nextPath != paths[2] {
		t.Error("got wrong next (third) image path")
	}
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
	nextImage, nextPath, err = il.Next()
	if err != nil {
		t.Errorf("could not get next (fourth) image. %s", err)
	}
	if nextPath != paths[3] {
		t.Error("got wrong next (fourth) image path")
	}
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
	nextImage, nextPath, err = il.Next()
	if err != nil {
		t.Errorf("could not get next (fifth) image. %s", err)
	}
	if nextPath != paths[4] {
		t.Error("got wrong next (fifth) image path")
	}
	if !compareImages(nextImage, expectImages[4]) {
		t.Error("got wrong next (fifth) image")
	}

	// Getting next image beyond length of list should return the last image
	nextImage, nextPath, err = il.Next()
	if err != nil {
		t.Errorf("could not get next (last) image after index is max. %s", err)
	}
	if nextPath != paths[4] {
		t.Error("got wrong next (fifth) image path after index is max")
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
	prevImage, prevPath, err := il.Previous()
	if err != nil {
		t.Errorf("could not get previous (fourth) image. %s", err)
	}
	if prevPath != paths[3] {
		t.Error("got wrong next (fourth) image path")
	}
	if !compareImages(prevImage, expectImages[3]) {
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
	prevImage, prevPath, err = il.Previous()
	if err != nil {
		t.Errorf("could not get previous (third) image. %s", err)
	}
	if prevPath != paths[2] {
		t.Error("got wrong next (third) image path")
	}
	if !compareImages(prevImage, expectImages[2]) {
		t.Error("got wrong previous (third) image")
	}
	// Get second image
	prevImage, prevPath, err = il.Previous()
	if err != nil {
		t.Errorf("could not get previous (second) image. %s", err)
	}
	if prevPath != paths[1] {
		t.Error("got wrong next (second) image path")
	}
	if !compareImages(prevImage, expectImages[1]) {
		t.Error("got wrong previous (second) image")
	}
	// Get first image
	prevImage, prevPath, err = il.Previous()
	if err != nil {
		t.Errorf("could not get previous (first) image. %s", err)
	}
	if prevPath != paths[0] {
		t.Error("got wrong next (first) image path")
	}
	if !compareImages(prevImage, expectImages[0]) {
		t.Error("got wrong previous (first) image")
	}
	// Getting previous image before zero should return the first image
	prevImage, prevPath, err = il.Previous()
	if err != nil {
		t.Errorf("could not get previous (first) image after index is zero. %s", err)
	}
	if prevPath != paths[0] {
		t.Error("got wrong next (first) image path after index is zero")
	}
	if !compareImages(prevImage, expectImages[0]) {
		t.Error("got wrong previous (first) image after index is zero")
	}

	// First image in image cache should be first image in list
	if !compareImages(il.imageCache[0], expectImages[0]) {
		t.Error("first image in image cache should be first image in path list")
	}
	// Last image in image cache should be third image in list
	if !compareImages(il.imageCache[2], expectImages[2]) {
		t.Error("last image in image cache should be third image in path list")
	}

	// Index should not be less than zero
	if il.index != 0 {
		t.Error("index should be zero")
	}
}

func TestLoadImage(t *testing.T) {
	// Load non-existant file
	_, err := loadImage("notAfile.jpg")
	if err == nil {
		t.Error("loading a non-existant file should return an error")
	}
	_, err = loadImage("./test_images/notAnImage.jpg")
	if err == nil {
		t.Error("loading a file that is not an image should return an error")
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
