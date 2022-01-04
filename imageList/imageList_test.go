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
	paths := []string{"./test_images/Obi1.jpg", "./test_images/Obi2.jpg"}

	imgList := make([]image.Image, 2)
	for i, path := range paths[:2] {
		imgFile, err := os.Open(path)
		if err != nil {
			t.Errorf("failed to open %s %s", path, err)
		}
		img, _, err := image.Decode(imgFile)
		if err != nil {
			t.Errorf("failed to decode %s %s", path, err)
		}
		imgList[i] = img
		imgFile.Close()
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

func TestNext(t *testing.T) {
	paths := []string{
		"./test_images/Obi1.jpg",
		"./test_images/Obi2.jpg",
		"./test_images/Jango1.jpg",
		"./test_images/Jango2.jpg",
	}
	expectImages := [4]image.Image{}
	for i, p := range paths {
		img, _ := loadImage(p)
		expectImages[i] = img
	}

	il, _ := New(paths)
	nextImage, err := il.Next()
	if err != nil {
		t.Errorf("could not get next image %s", err)
	}

	// nextImage should be the first image in the list
	if !compareImages(nextImage, expectImages[0]) {
		t.Error("got wrong next image")
	}

	// Index should be 1
	if il.index != 1 {
		t.Error("index was not iterated. index =", il.index)
	}

	// Image cache should contain the first three images in the list
	for i, img := range expectImages[:3] {
		if !compareImages(img, il.imageCache[i]) {
			t.Error("image cache is not correct at index", i)
		}
	}
	il.Next()
	il.Next()

	// First image in the image cache should be the second image in the list
	if !compareImages(il.imageCache[0], expectImages[1]) {
		t.Error("first image in image cache should be second image in path list")
	}

	// Last image in the image cache should be the last image in the list
	if !compareImages(il.imageCache[len(il.imageCache)-1], expectImages[len(expectImages)-1]) {
		t.Error("last image in the image cache should be last image in the path list")
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
