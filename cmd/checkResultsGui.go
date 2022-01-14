package cmd

import (
	"image"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	imagelist "github.com/mike-lloyd03/dedugo/imageList"
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
	refIL         imagelist.ImageList
	dupeIL        imagelist.ImageList
)

func showGui() {
	results = readResultsFile(results_path)

	w := a.NewWindow("dedugo")
	w.CenterOnScreen()

	var refImages, dupeImages []string
	for _, p := range results.ImagePairs {
		refImages = append(refImages, p.RefImage)
		dupeImages = append(dupeImages, p.DupeImage)
	}

	var err error
	refIL, err = imagelist.New(refImages)
	if err != nil {
		log.Fatal("error creating reference image list.", err)
	}
	dupeIL, err = imagelist.New(dupeImages)
	if err != nil {
		log.Fatal("error creating duplicate image list.", err)
	}
	refImgFile, refPath, err := refIL.Next()
	if err != nil {
		log.Fatal("failed to get initial referance image.", err)
	}
	dupeImgFile, dupePath, err := dupeIL.Next()
	if err != nil {
		log.Fatal("failed to get initial duplicate image.", err)
	}

	refLabel := widget.NewLabelWithStyle("Reference Image", textCentered, bold)
	refImagePath = widget.NewLabelWithStyle(refPath, textCentered, monospaced)
	refImage = canvas.NewImageFromImage(refImgFile)
	refImage.SetMinSize(fyne.NewSize(imgWidth, imgHeight))
	refImage.FillMode = canvas.ImageFillContain

	dupeLabel := widget.NewLabelWithStyle("Duplicate Image", textCentered, bold)
	dupeImagePath = widget.NewLabelWithStyle(dupePath, textCentered, monospaced)
	dupeImage = canvas.NewImageFromImage(dupeImgFile)
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
		p := results.ImagePairs[results.StartIdx]
		p.Confirmed = true
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
	var refImageFile image.Image
	var refPath string
	var dupeImageFile image.Image
	var dupePath string
	var err error

	if direction == Next {
		refImageFile, refPath, err = refIL.Next()
		if err != nil {
			log.Fatal("failed to get next referance image.", err)
		}
		dupeImageFile, dupePath, err = dupeIL.Next()
		if err != nil {
			log.Fatal("failed to get next duplicate image.", err)
		}
	} else if direction == Previous {
		refImageFile, refPath, err = refIL.Previous()
		if err != nil {
			log.Fatal("failed to get previous referance image.", err)
		}
		dupeImageFile, dupePath, err = dupeIL.Previous()
		if err != nil {
			log.Fatal("failed to get previous duplicate image.", err)
		}
	}
	refImage.Image = refImageFile
	refImagePath.Text = refPath
	dupeImage.Image = dupeImageFile
	dupeImagePath.Text = dupePath
	refImage.Refresh()
	refImagePath.Refresh()
	dupeImage.Refresh()
	dupeImagePath.Refresh()
}
