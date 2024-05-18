package gui

import (
	"CUBUS-core/shared"
	"CUBUS-core/translation"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/google/uuid"
)

func setupDialog(window fyne.Window) shared.CubeConfigType { // TODO: use the shared dialog tree to create the setup dialog form and return the cube config
	T := translation.T

	cubeConfig := shared.CubeConfigType{
		Id:        uuid.New().String(),
		CubeType:  shared.CubeTypes["generic-worker"],
		PublicKey: nil,
	}

	cubeTypeLabels := make([]string, len(shared.CubeTypes))
	counter := 0
	for _, cubeType := range shared.CubeTypes {
		cubeTypeLabels[counter] = cubeType.Label
		counter++
	}

	locationItem := widget.NewSelect([]string{T("Local"), T("Remote")}, func(string) {})
	typeItem := widget.NewSelect(cubeTypeLabels, func(string) {})
	publicKeyItem := widget.NewMultiLineEntry()

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: T("Location"), Widget: locationItem},
			{Text: T("Type"), Widget: typeItem},
			{Text: T("Public Key"), Widget: publicKeyItem},
		},
		OnSubmit: func() {
		},
	}

	dialog.NewCustom(T("Setup"), "OK", form, window).Show()

	return cubeConfig
}
