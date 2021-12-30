package gui

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
	isDuplicate bool
	a           = app.New()
)

func ShowGui(refImagePath, dupeImagePath string) bool {
	w := a.NewWindow("dedugo")
	w.CenterOnScreen()

	refLabel := widget.NewLabel("Reference Image")
	refImage := canvas.NewImageFromFile(refImagePath)
	refImage.SetMinSize(fyne.NewSize(imgWidth, imgHeight))
	refImage.FillMode = canvas.ImageFillContain

	dupeLabel := widget.NewLabel("Duplicate Image")
	dupeImage := canvas.NewImageFromFile(dupeImagePath)
	dupeImage.SetMinSize(fyne.NewSize(imgWidth, imgHeight))
	dupeImage.FillMode = canvas.ImageFillContain

	yesButton := widget.NewButton("Yes", yesFunc)
	noButton := widget.NewButton("No", noFunc)

	refImgCont := container.NewVBox(refLabel, refImage)
	dupeImgCont := container.NewVBox(dupeLabel, dupeImage)
	imgCont := container.NewHBox(refImgCont, dupeImgCont)
	buttonCont := container.NewHBox(layout.NewSpacer(), noButton, yesButton, layout.NewSpacer())
	mainCont := container.NewVBox(imgCont, buttonCont)
	w.SetContent(mainCont)

	w.ShowAndRun()
	return isDuplicate
}

func yesFunc() {
	isDuplicate = true
	a.Quit()
}

func noFunc() {
	isDuplicate = false
	a.Quit()
}
