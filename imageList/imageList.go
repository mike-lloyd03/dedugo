package imagelist

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"io/ioutil"
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

	imgCache := make([]image.Image, 2)
	for i, path := range paths[:2] {
		img, err := loadImage(path)
		if err != nil {
			return ImageList{}, errors.New(fmt.Sprint("could not load image", path))
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
	if il.index == len(il.paths) {
		return nil, errors.New("end of path array")
	}
	nextImage := il.imageCache[il.index]
	il.index++
	appendImage, _ := loadImage(il.paths[il.index+1])
	fmt.Println("appending", il.paths[il.index+1])
	il.imageCache = append(il.imageCache[len(il.imageCache)-2:len(il.imageCache)], appendImage)
	return nextImage, nil
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
