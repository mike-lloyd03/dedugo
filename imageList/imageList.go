package imagelist

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"io/ioutil"
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

func New(paths []string) (ImageList, error) {
	if len(paths) < 2 {
		return ImageList{}, errors.New("paths argument must contain at least 2 items")
	}

	imgCache := make([]image.Image, 3)
	for i := 0; i < 3; i++ {
		img, err := loadImage(paths[i])
		if err != nil {
			return ImageList{}, errors.New(fmt.Sprint("could not load image", paths[i]))
		}
		imgCache[i] = img
	}

	return ImageList{
		index:      0,
		paths:      paths,
		imageCache: imgCache,
	}, nil
}

func (il *ImageList) Next() (image.Image, error) {
	// TODO: If list is less than three images
	w.Wait()
	if il.index < 2 {
		nextImage := il.imageCache[il.index]
		il.index++
		return nextImage, nil
	} else if 2 <= il.index && il.index < len(il.paths)-1 {
		nextImage := il.imageCache[2]
		il.index++
		w.Add(1)
		go il.appendImage()
		return nextImage, nil
	} else if il.index == len(il.paths)-1 {
		nextImage := il.imageCache[2]
		return nextImage, nil
	}
	return nil, errors.New(fmt.Sprint("image list index out of range: ", il.index))
}

func (il *ImageList) Previous() (image.Image, error) {
	// TODO: If list is less than three images
	w.Wait()
	if il.index == len(il.paths)-1 {
		prevImage := il.imageCache[1]
		il.index--
		return prevImage, nil
	} else if 1 < il.index && il.index < len(il.paths)-1 {
		prevImage := il.imageCache[0]
		il.index--
		w.Add(1)
		go il.prependImage()
		return prevImage, nil
	} else if il.index == 1 {
		prevImage := il.imageCache[0]
		il.index--
		return prevImage, nil
	} else if il.index == 0 {
		prevImage := il.imageCache[0]
		return prevImage, nil
	}
	return nil, errors.New(fmt.Sprint("image list index out of range: ", il.index))
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
