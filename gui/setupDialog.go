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

func setupDialog(window fyne.Window) types.CubeConfig { // TODO: use onSubmit of form to get form data and return it
	T := translation.T

	cubeConfig := types.CubeConfig{
		Id:        uuid.New().String(),
		CubeType:  types.CubeTypes.GenericWorker,
		PublicKey: nil,
	}

	cubeSetupForm := forms.GetCubeSetupForm()
	box := container.New(layout.NewVBoxLayout())

	formPopup := dialog.NewCustomWithoutButtons(T("Setup"), box, window)
	formPopup.Resize(fyne.NewSize(700, 400))
	forms.FormToFyneForm(cubeSetupForm, box, formPopup, window)
	formPopup.Show()

	return cubeConfig
}
