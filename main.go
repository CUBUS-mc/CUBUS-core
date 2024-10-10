package main

import (
	"CUBUS-core/gui"
	"CUBUS-core/orchestrator/server"
	"CUBUS-core/shared"
	"CUBUS-core/shared/translation"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	T := translation.T
	cubusApp := app.NewWithID("com.virusrpi.cubus-core")
	icon, _ := fyne.LoadResourceFromURLString("https://raw.githubusercontent.com/CUBUS-mc/CUBUS-core/master/assets/android.png")
	cubusApp.SetIcon(icon)
	defaultValues := shared.NewDefaults()

	ui := shared.UiType(cubusApp.Preferences().StringWithFallback("ui", string(defaultValues.UI)))
	cubusApp.Preferences().SetString("ui", string(ui))
	language := cubusApp.Preferences().StringWithFallback("language", defaultValues.Language)
	cubusApp.Preferences().SetString("language", language)
	translation.ChangeLanguage(language)

	localOrchestratorServer := server.NewServer(":25560")
	localOrchestratorServer.Start()

	switch ui {
	case shared.TUI:
		println(T("TUI not implemented yet")) // Implement TUI with https://charm.sh/
	case shared.GUI:
		gui.NewGui(cubusApp, defaultValues).Run()
	case shared.NONE: // Do nothing
	}
}
