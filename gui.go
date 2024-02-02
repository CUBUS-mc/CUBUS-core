package main

import "fyne.io/fyne/v2/app"

func gui(configuration map[string]interface{}) {
	qubusApp := app.New()
	qubusWindow := qubusApp.NewWindow("QUBUS")

	qubusWindow.ShowAndRun()
}
