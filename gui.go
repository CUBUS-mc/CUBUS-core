package main

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

func selectCube(c *cube, infoContainerShape *canvas.Rectangle, pointerLine *canvas.Line, pointerTip *canvas.Circle) {
	if c.selected {
		canvas.NewPositionAnimation(
			infoContainerShape.Position(),
			fyne.NewPos(1400, 25),
			time.Second/4,
			func(pos fyne.Position) {
				infoContainerShape.Move(pos)
				infoContainerShape.Refresh()
			}).Start()
		pointerTip.Hide()
		pointerLine.Hide()
		c.selected = false
	} else {
		canvas.NewPositionAnimation(
			infoContainerShape.Position(),
			fyne.NewPos(1400-325, 25),
			time.Second/4,
			func(pos fyne.Position) {
				infoContainerShape.Move(pos)
				infoContainerShape.Refresh()
			}).Start()
		pointerTip.Move(fyne.NewPos(c.Position().X+c.size/2-5, c.Position().Y+c.size/2-5))
		pointerTip.Show()
		pointerLine.Resize(fyne.NewSize(0, 0))
		pointerLine.Move(fyne.NewPos(c.Position().X+c.size/2, c.Position().Y+c.size/2))
		pointerLine.Show()
		canvas.NewSizeAnimation(
			pointerLine.Size(),
			fyne.NewSize(1400-c.Position().X-325, 2),
			time.Second/4,
			func(size fyne.Size) {
				pointerLine.Resize(size)
				pointerLine.Refresh()
				if size.Width == 1400-c.Position().X-325 {
					c.selected = true
				}
			}).Start()
	}
}

func gui(qubusApp fyne.App, defaults *Defaults) {
	qubusWindow := qubusApp.NewWindow("QUBUS core")
	qubusWindow.Resize(fyne.NewSize(1400, 900))
	qubusWindow.CenterOnScreen()
	qubusWindow.SetIcon(qubusApp.Icon())
	qubusWindow.SetFixedSize(true)

	windowMenu := fyne.NewMainMenu(
		fyne.NewMenu("File",
			fyne.NewMenuItem("Create a new Qube", func() {
				println("Create")
			}),
			fyne.NewMenuItem("Export config", func() {
				println("Export config")
			}),
			fyne.NewMenuItem("Import config", func() {
				println("Import config")
			}),
			fyne.NewMenuItem("Settings", func() {
				println("Settings")
			}),
		),
	)
	qubusWindow.SetMainMenu(windowMenu)

	infoContainerShape := canvas.NewRectangle(color.White)
	infoContainerShape.Resize(fyne.NewSize(300, 825))
	infoContainerShape.Move(fyne.NewPos(1400, 25))
	infoContainerShape.CornerRadius = 12

	pointerLine := canvas.NewLine(color.White)
	pointerLine.StrokeWidth = 2
	pointerLine.Hide()
	pointerTip := canvas.NewCircle(color.White)
	pointerTip.Resize(fyne.NewSize(10, 10))
	pointerTip.Hide()

	cubeImage := newCube(defaults.CubeAssetURL, func(c *cube) { selectCube(c, infoContainerShape, pointerLine, pointerTip) })

	mainContainer := container.NewWithoutLayout(cubeImage, infoContainerShape, pointerLine, pointerTip)
	qubusWindow.SetContent(mainContainer)

	go func() {
		for {
			mainContainer.Refresh()
			pointerTip.Refresh()
			if cubeImage.selected {
				pointerTip.Move(fyne.NewPos(cubeImage.Position().X+cubeImage.size/2-5, cubeImage.Position().Y+cubeImage.size/2-5))
				pointerLine.Move(fyne.NewPos(cubeImage.Position().X+cubeImage.size/2, cubeImage.Position().Y+cubeImage.size/2))
				pointerLine.Resize(fyne.NewSize(1400-cubeImage.Position().X-325, 2))
			}
		}
	}()

	qubusWindow.ShowAndRun()
}
