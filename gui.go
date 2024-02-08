package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

func gui(qubusApp fyne.App, defaults *Defaults) {
	qubusWindow := qubusApp.NewWindow("QUBUS core")
	qubusWindow.Resize(fyne.NewSize(800, 600))
	qubusWindow.CenterOnScreen()
	qubusWindow.SetIcon(qubusApp.Icon())

	windowMenu := fyne.NewMainMenu(
		fyne.NewMenu("File",
			fyne.NewMenuItem("Create a new Qube", func() {
				println("Create")
			}),
			fyne.NewMenuItem("Export config", func() {
				println("Export config")
			}),
		),
	)
	qubusWindow.SetMainMenu(windowMenu)

	cubeImageResource, _ := fyne.LoadResourceFromURLString(defaults.CubeAssetURL)
	cubeImage := canvas.NewImageFromResource(cubeImageResource)
	cubeImage.Move(fyne.NewPos(400, 300))
	cubeImage.Resize(fyne.NewSize(150, 150))

	mainContainer := container.NewWithoutLayout(cubeImage)
	qubusWindow.SetContent(mainContainer)

	qubusWindow.ShowAndRun()
}
