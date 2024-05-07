package main

import (
	"CUBUS/gui"
	"CUBUS/shared"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	cubusApp := app.NewWithID("com.virusrpi.cubus")
	icon, _ := fyne.LoadResourceFromURLString("https://raw.githubusercontent.com/CUBUS-mc/CUBUS-core/master/assets/android.png")
	cubusApp.SetIcon(icon)
	defaultValues := shared.NewDefaults()

	ui := shared.UiType(cubusApp.Preferences().StringWithFallback("ui", string(defaultValues.UI)))
	cubusApp.Preferences().SetString("ui", string(ui))

	switch ui {
	case shared.CLI:
		println("CLI not implemented yet")
	case shared.GUI:
		gui.Gui(cubusApp, defaultValues)
	case shared.API:
		println("API not implemented yet")
	}
}
