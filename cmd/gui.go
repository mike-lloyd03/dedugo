package cmd

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

const (
	imgWidth  = 640
	imgHeight = 480
)

var (
	textCentered = fyne.TextAlignCenter
	monospaced   = fyne.TextStyle{Monospace: true}
	bold         = fyne.TextStyle{Bold: true}
	a            = app.New()
	refImage     *canvas.Image
	dupeImage    *canvas.Image
	results      Results
)

func showGui() {
	results = readResultsFile(results_path)
	p := results.ImagePairs[results.StartIdx]

	w := a.NewWindow("dedugo")
	w.CenterOnScreen()

	refLabel := widget.NewLabelWithStyle("Reference Image", textCentered, bold)
	refImagePath := widget.NewLabelWithStyle(p.RefImage, textCentered, monospaced)
	refImage = canvas.NewImageFromFile(p.RefImage)
	refImage.SetMinSize(fyne.NewSize(imgWidth, imgHeight))
	refImage.FillMode = canvas.ImageFillContain

	dupeLabel := widget.NewLabelWithStyle("Duplicate Image", textCentered, bold)
	dupeImagePath := widget.NewLabelWithStyle(p.DupeImage, textCentered, monospaced)
	dupeImage = canvas.NewImageFromFile(p.DupeImage)
	dupeImage.SetMinSize(fyne.NewSize(imgWidth, imgHeight))
	dupeImage.FillMode = canvas.ImageFillContain

	refImgCont := container.NewVBox(refLabel, refImagePath, refImage)
	dupeImgCont := container.NewVBox(dupeLabel, dupeImagePath, dupeImage)
	imgCont := container.NewHBox(refImgCont, dupeImgCont)

	confirmButton := widget.NewButton("Confirm Duplicate", confirmDuplicate())
	nextButton := widget.NewButton("Next", nextPair())
	prevButton := widget.NewButton("Previous", prevPair())

	buttonCont := container.NewHBox(layout.NewSpacer(), prevButton, nextButton, confirmButton, layout.NewSpacer())
	mainCont := container.NewVBox(imgCont, buttonCont)

	w.SetContent(mainCont)

	w.ShowAndRun()
}

func confirmDuplicate() func() {
	return func() {
		p := results.ImagePairs[results.StartIdx]
		p.Confirmed = true
		results.StartIdx++
		WriteResultsFile(results, results_path)

		refreshImages(results)
	}
}

func nextPair() func() {
	return func() {
		results.StartIdx++
		WriteResultsFile(results, results_path)

		refreshImages(results)
	}
}

func prevPair() func() {
	return func() {
		results.StartIdx--
		WriteResultsFile(results, results_path)

		refreshImages(results)
	}
}

func refreshImages(results Results) {
	p := results.ImagePairs[results.StartIdx]
	refImage.File = p.RefImage
	dupeImage.File = p.DupeImage
	refImage.Refresh()
	dupeImage.Refresh()
}
