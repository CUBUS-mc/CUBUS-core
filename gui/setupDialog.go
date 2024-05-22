package gui

import (
	"CUBUS-core/shared/forms"
	"CUBUS-core/shared/translation"
	"CUBUS-core/shared/types"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"github.com/google/uuid"
)

func setupDialog(window fyne.Window) types.CubeConfig { // TODO: add a on change listener that puts the values into the cubeConfig
	T := translation.T

	cubeConfig := types.CubeConfig{
		Id:        uuid.New().String(),
		CubeType:  types.CubeTypes.GenericWorker,
		PublicKey: nil,
	}

	cubeSetupForm := forms.GetCubeSetupForm()
	box := container.New(layout.NewVBoxLayout())
	forms.FormToFyneForm(cubeSetupForm, box)

	dialog.NewCustom(T("Setup"), "OK", box, window).Show()

	return cubeConfig
}
