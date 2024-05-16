package main

import (
	"CUBUS-core/gui"
	"CUBUS-core/shared"
	"CUBUS-core/translation"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	T := translation.T
	cubusApp := app.NewWithID("com.virusrpi.cubus-core")
	icon, _ := fyne.LoadResourceFromURLString("https://raw.githubusercontent.com/CUBUS-mc/CUBUS-core/master/assets/android.png")
	cubusApp.SetIcon(icon)
	defaultValues := shared.NewDefaults()

	ui := shared.CLI //shared.UiType(cubusApp.Preferences().StringWithFallback("ui", string(defaultValues.UI)))
	cubusApp.Preferences().SetString("ui", string(ui))

	switch ui {
	case shared.CLI:
		println(T("CLI not implemented yet"))
	case shared.GUI:
		gui.Gui(cubusApp, defaultValues)
	case shared.API:
		println(T("API not implemented yet"))
	}
}
