package main

import "fyne.io/fyne/v2/app"

func gui(configuration map[string]interface{}) {
	println(configuration["cube_type"])

	qubusApp := app.New()
	qubusWindow := qubusApp.NewWindow("QUBUS")

	qubusWindow.ShowAndRun()
}
