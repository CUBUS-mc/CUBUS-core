package gui

import (
	"CUBUS-core/shared/forms"
	"CUBUS-core/shared/translation"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
)

func settingsDialog(cubusApp fyne.App) {
	currentLanguage := cubusApp.Preferences().String("language")
	settingsForm := forms.GetUiSettingsForm(currentLanguage)
	box := container.New(layout.NewVBoxLayout())

	formSubmitCallback := func(values map[string]string) {
		cubusApp.Preferences().SetString("language", values["language"])
		translation.ChangeLanguage(cubusApp.Preferences().String("language"))
	}

	formPopup := dialog.NewCustomWithoutButtons("Settings", box, cubusWindow)
	formPopup.Resize(fyne.NewSize(WindowWidth()*0.5, WindowHeight()*0.5))
	forms.FormToFyneForm(
		settingsForm,
		box,
		formPopup,
		cubusWindow,
		formSubmitCallback,
	)
	formPopup.Show()
}
