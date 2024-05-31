package gui

import (
	"CUBUS-core/shared"
	"CUBUS-core/shared/forms"
	"CUBUS-core/shared/translation"
	"CUBUS-core/shared/types"
	"CUBUS-core/shared/types/gui"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/google/uuid"
)

func setupDialog(
	window fyne.Window,
	cubeConfigs *[]map[string]interface{},
	cubeStrings *[]string,
	cubusApp fyne.App,
	defaults *shared.Defaults,
	cubeContainerObject *gui.CubeContainer,
	infoContainerShape *canvas.Rectangle,
	pointerLine *canvas.Line,
	pointerTip *canvas.Circle,
	infoContainerText *widget.RichText,
) {
	T := translation.T

	cubeConfig := types.CubeConfig{
		Id:        uuid.New().String(),
		CubeType:  types.CubeTypes.GenericWorker,
		CubeName:  "",
		PublicKey: nil,
	}

	cubeSetupForm := forms.GetCubeSetupForm()
	box := container.New(layout.NewVBoxLayout())

	formSubmitCallback := func(
		values map[string]string,
		cubeConfigs *[]map[string]interface{},
		cubeStrings *[]string,
		cubusApp fyne.App,
		defaults shared.Defaults,
		cubeContainerObject *gui.CubeContainer,
		infoContainerShape *canvas.Rectangle,
		pointerLine canvas.Line,
		pointerTip canvas.Circle,
		infoContainerText *widget.RichText,
	) {
		for key, value := range values {
			switch key {
			case "cubeName":
				cubeConfig.CubeName = value
				break
			case "cubeType":
				cubeConfig.CubeType = types.CubeType{Value: value}
				break
			}
		}
		*cubeConfigs = append(
			*cubeConfigs,
			map[string]interface{}{
				"id":   cubeConfig.Id,
				"name": cubeConfig.CubeName,
				"type": cubeConfig.CubeType.Value,
			},
		)
		*cubeStrings = []string{}
		for _, cubeConfig := range *cubeConfigs {
			*cubeStrings = append(*cubeStrings, shared.ObjectToJsonString(cubeConfig))
		}
		cubusApp.Preferences().SetStringList("cubes", *cubeStrings)
		cubeContainerObject.AddCube(defaults.CubeAssetURL, func(c *gui.Cube) { selectCube(c, infoContainerShape, &pointerLine, &pointerTip, infoContainerText) }, cubeConfig.Id, cubeConfig)
		cubeContainerObject.CenterCubes()
	}

	formPopup := dialog.NewCustomWithoutButtons(T("Setup"), box, window)
	formPopup.Resize(fyne.NewSize(700, 400))
	forms.FormToFyneForm(
		cubeSetupForm,
		box,
		formPopup,
		window,
		formSubmitCallback,
		cubeConfigs,
		cubeStrings,
		cubusApp,
		defaults,
		cubeContainerObject,
		infoContainerShape,
		pointerLine,
		pointerTip,
		infoContainerText,
	)
	formPopup.Show()
}
