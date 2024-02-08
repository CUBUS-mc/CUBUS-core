package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	cubusApp := app.NewWithID("com.virusrpi.cubus")
	icon, _ := fyne.LoadResourceFromURLString("https://raw.githubusercontent.com/CUBUS-mc/CUBUS-core/master/assets/android.png")
	cubusApp.SetIcon(icon)
	defaultValues := NewDefaults()

	ui := UiType(cubusApp.Preferences().StringWithFallback("ui", string(defaultValues.UI)))
	cubusApp.Preferences().SetString("ui", string(ui))

	switch ui {
	case CLI:
		// CLI
	case GUI:
		gui(cubusApp, defaultValues)
	case API:
		// API
	}
}
