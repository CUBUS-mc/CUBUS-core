package main

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

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

	cubeImage := newCube(defaults.CubeAssetURL)

	infoContainerShape := canvas.NewRectangle(color.White)
	infoContainerShape.Resize(fyne.NewSize(300, 825))
	infoContainerShape.Move(fyne.NewPos(1400, 25))
	infoContainerShape.CornerRadius = 12

	mainContainer := container.NewWithoutLayout(cubeImage, infoContainerShape)
	qubusWindow.SetContent(mainContainer)

	canvas.NewPositionAnimation(
		infoContainerShape.Position(),
		fyne.NewPos(1400-325, 25),
		time.Second,
		func(pos fyne.Position) {
			infoContainerShape.Move(pos)
			infoContainerShape.Refresh()
		}).Start()

	go func() {
		for {
			mainContainer.Refresh()
		}
	}()

	qubusWindow.ShowAndRun()
}
