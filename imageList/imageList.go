package imagelist

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"io/ioutil"
	"math"
	"sync"
)

var (
	m sync.Mutex
	w sync.WaitGroup
)

type ImageList struct {
	index      int
	paths      []string
	imageCache []image.Image
}

// New creates a new ImageList type and initializes the image cache.
func New(paths []string, index ...int) (ImageList, error) {
	var i int
	if len(index) == 0 {
		i = 0
	} else {
		i = index[0]
	}

	if len(paths) < 3 {
		return ImageList{}, errors.New("paths argument must contain at least 2 items")
	}

	il := ImageList{
		index: i,
		paths: paths,
	}

	err := il.initCache(i)
	if err != nil {
		return ImageList{}, errors.New(fmt.Sprintf("failed to initialize image cache. %s", err))
	}

	return il, nil
}

// Next returns the next image in the cache and it's path.
func (il *ImageList) Next() (image.Image, string, error) {
	// TODO: If list is less than three images
	w.Wait()
	nextPath := il.paths[il.index]
	if il.index < 2 {
		nextImage := il.imageCache[il.index]
		il.index++
		return nextImage, nextPath, nil
	} else if 2 <= il.index && il.index < len(il.paths)-1 {
		nextImage := il.imageCache[2]
		il.index++
		w.Add(1)
		go il.appendImage()
		return nextImage, nextPath, nil
	} else if il.index == len(il.paths)-1 {
		nextImage := il.imageCache[2]
		return nextImage, nextPath, nil
	}
	return nil, "", errors.New(fmt.Sprint("image list index out of range: ", il.index))
}

// Previous returns the previous image in the cache and it's path.
func (il *ImageList) Previous() (image.Image, string, error) {
	// TODO: If list is less than three images
	prevPath := il.paths[int(math.Max(0, float64(il.index-1)))]
	w.Wait()
	if il.index == len(il.paths)-1 {
		prevImage := il.imageCache[1]
		il.index--
		return prevImage, prevPath, nil
	} else if 1 < il.index && il.index < len(il.paths)-1 {
		prevImage := il.imageCache[0]
		il.index--
		w.Add(1)
		go il.prependImage()
		return prevImage, prevPath, nil
	} else if il.index == 1 {
		prevImage := il.imageCache[0]
		il.index--
		return prevImage, prevPath, nil
	} else if il.index == 0 {
		prevImage := il.imageCache[0]
		return prevImage, prevPath, nil
	}
	return nil, "", errors.New(fmt.Sprint("image list index out of range: ", il.index))
}

func (il *ImageList) appendImage() {
	defer w.Done()
	aImg, _ := loadImage(il.paths[il.index])
	m.Lock()
	il.imageCache = append(il.imageCache[1:], aImg)
	m.Unlock()
}

func (il *ImageList) prependImage() {
	defer w.Done()
	pImg, _ := loadImage(il.paths[il.index-1])
	m.Lock()
	il.imageCache = append([]image.Image{pImg}, il.imageCache[:2]...)
	m.Unlock()
}

func (il *ImageList) initCache(index int) error {
	imgCache := make([]image.Image, 3)
	for i := index; i < index+3; i++ {
		img, err := loadImage(il.paths[i])
		if err != nil {
			return errors.New(fmt.Sprintf("could not load image. %s", err))
		}
		imgCache[i] = img
	}
	il.imageCache = imgCache
	return nil
}

func (il *ImageList) goToIndex(index int) {
}

func loadImage(path string) (image.Image, error) {
	imgBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.New("image could not be opened")
	}
	imgReader := bytes.NewReader(imgBytes)
	img, _, err := image.Decode(imgReader)
	if err != nil {
		return nil, errors.New("image could not be decoded")
	}
	return img, nil
}
