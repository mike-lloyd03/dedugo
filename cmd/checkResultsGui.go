package cmd

import (
	"bytes"
	"errors"
	"image"
	"io/ioutil"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type ImageReader struct {
	Path   string
	Reader image.Image
}

type Direction int

const (
	imgWidth            = 640
	imgHeight           = 480
	Previous  Direction = iota
	Current
	Next
)

var (
	textCentered  = fyne.TextAlignCenter
	monospaced    = fyne.TextStyle{Monospace: true}
	bold          = fyne.TextStyle{Bold: true}
	a             = app.New()
	results       Results
	refImage      *canvas.Image
	dupeImage     *canvas.Image
	refImagePath  *widget.Label
	dupeImagePath *widget.Label
	nextButton    *widget.Button
	prevButton    *widget.Button
	prevRefImage  image.Image
	currRefImage  image.Image
	nextRefImage  image.Image
	prevDupeImage image.Image
	currDupeImage image.Image
	nextDupeImage image.Image
)

func showGui() {
	results = readResultsFile(results_path)
	p := results.ImagePairs[results.StartIdx]

	w := a.NewWindow("dedugo")
	w.CenterOnScreen()

	initImages()

	refLabel := widget.NewLabelWithStyle("Reference Image", textCentered, bold)
	refImagePath = widget.NewLabelWithStyle(p.RefImage, textCentered, monospaced)
	refImage = canvas.NewImageFromImage(currRefImage)
	refImage.SetMinSize(fyne.NewSize(imgWidth, imgHeight))
	refImage.FillMode = canvas.ImageFillContain

	dupeLabel := widget.NewLabelWithStyle("Duplicate Image", textCentered, bold)
	dupeImagePath = widget.NewLabelWithStyle(p.DupeImage, textCentered, monospaced)
	dupeImage = canvas.NewImageFromImage(currDupeImage)
	dupeImage.SetMinSize(fyne.NewSize(imgWidth, imgHeight))
	dupeImage.FillMode = canvas.ImageFillContain

	refImgCont := container.NewVBox(refLabel, refImagePath, refImage)
	dupeImgCont := container.NewVBox(dupeLabel, dupeImagePath, dupeImage)
	imgCont := container.NewHBox(refImgCont, dupeImgCont)

	confirmButton := widget.NewButton("Confirm Duplicate", confirmDuplicate())
	nextButton = widget.NewButton("Next", nextPair())
	prevButton = widget.NewButton("Previous", prevPair())

	buttonCont := container.NewHBox(layout.NewSpacer(), prevButton, nextButton, confirmButton, layout.NewSpacer())
	mainCont := container.NewVBox(imgCont, buttonCont)

	w.SetContent(mainCont)

	w.ShowAndRun()
}

func confirmDuplicate() func() {
	return func() {
		results.ImagePairs[results.StartIdx].Confirmed = true
		go WriteResultsFile(results, results_path)
		nextPair()
	}
}

func nextPair() func() {
	return func() {
		if results.StartIdx <= len(results.ImagePairs) {
			results.StartIdx++
			go WriteResultsFile(results, results_path)
			refreshImages(Next)
		} else {
			nextButton.Disable()
		}
	}
}

func prevPair() func() {
	return func() {
		results.StartIdx--
		go WriteResultsFile(results, results_path)

		refreshImages(Previous)
	}
}

func refreshImages(direction Direction) {
	if direction == Next {
		prevRefImage = currRefImage
		prevDupeImage = currDupeImage
		currRefImage = nextRefImage
		currDupeImage = nextDupeImage
		nextRefImage, nextDupeImage = loadImage(Next)
	} else if direction == Previous {
		nextRefImage = currRefImage
		nextDupeImage = currDupeImage
		currRefImage = prevRefImage
		currDupeImage = prevDupeImage
		prevRefImage, prevDupeImage = loadImage(Previous)
	}
	refImage.Image = currRefImage
	refImagePath.Text = results.ImagePairs[results.StartIdx].RefImage
	dupeImage.Image = currDupeImage
	dupeImagePath.Text = results.ImagePairs[results.StartIdx].DupeImage
	go refImage.Refresh()
	go refImagePath.Refresh()
	go dupeImage.Refresh()
	go dupeImagePath.Refresh()
}

func initImages() {
	i := results.StartIdx
	if i == 0 {
		prevRefImage, prevDupeImage = nil, nil
		currRefImage, currDupeImage = loadImage(Current)
		nextRefImage, nextDupeImage = loadImage(Next)
	} else if i == len(results.ImagePairs) {
		prevRefImage, prevDupeImage = loadImage(Previous)
		currRefImage, currDupeImage = loadImage(Current)
		nextRefImage, nextDupeImage = nil, nil
	} else {
		prevRefImage, prevDupeImage = loadImage(Previous)
		currRefImage, currDupeImage = loadImage(Current)
		nextRefImage, nextDupeImage = loadImage(Next)
	}
}

func loadImage(direction Direction) (image.Image, image.Image) {
	var i int

	switch direction {
	case Next:
		i = results.StartIdx + 1
	case Previous:
		i = results.StartIdx - 1
	default:
		i = results.StartIdx
	}

	refImageImage, err := openAndDecodeImage(results.ImagePairs[i].RefImage)
	if err != nil {
		log.Fatal("Reference", err)
	}
	dupeImageImage, err := openAndDecodeImage(results.ImagePairs[i].DupeImage)
	if err != nil {
		log.Fatal("Duplicate", err)
	}
	return refImageImage, dupeImageImage
}

func openAndDecodeImage(path string) (image.Image, error) {
	imageBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.New("Image could not be opened.")
	}
	imageReader := bytes.NewReader(imageBytes)
	image, _, err := image.Decode(imageReader)
	if err != nil {
		return nil, errors.New("Image could not be decoded.")
	}
	return image, nil
}
