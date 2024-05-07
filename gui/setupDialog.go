package gui

import (
	"CUBUS/shared"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/google/uuid"
)

func setupDialog(window fyne.Window) shared.CubeConfigType { // TODO: use the shared dialog tree to create the setup dialog form and return the cube config
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

	locationItem := widget.NewSelect([]string{"Local", "Remote"}, func(string) {})
	typeItem := widget.NewSelect(cubeTypeLabels, func(string) {})
	publicKeyItem := widget.NewMultiLineEntry()

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Location", Widget: locationItem},
			{Text: "Type", Widget: typeItem},
			{Text: "Public Key", Widget: publicKeyItem},
		},
		OnSubmit: func() {
		},
	}

	dialog.NewCustom("Setup", "OK", form, window).Show()

	return cubeConfig
}
