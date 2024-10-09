package gui

import (
	"CUBUS-core/orchestrator/client"
	"CUBUS-core/shared"
	"CUBUS-core/shared/forms"
	"CUBUS-core/shared/translation"
	"CUBUS-core/shared/types"
	"CUBUS-core/shared/types/gui"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"github.com/google/uuid"
)

func setupDialog(
	window fyne.Window,
	cubeConfigs *[]types.CubeConfig,
	defaults *shared.Defaults,
	cubeContainerObject *gui.CubeContainer,
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

		orchestratorClient := client.NewClient()
		var serverUrl string
		if values["cubeLocation"] == "local" {
			serverUrl = "http://localhost:25560"
		} else {
			serverUrl = values["remoteUrl"]
		}
		err := orchestratorClient.CreateNewCube(serverUrl, cubeConfig)
		if err != nil {
			dialog.ShowError(err, window)
			return
		}

		*cubeConfigs = append(*cubeConfigs, cubeConfig)
		cubeContainerObject.AddCube(defaults.CubeAssetURL, cubeConfig.Id, cubeConfig)
		cubeContainerObject.CenterCubes()
	}

	formPopup := dialog.NewCustomWithoutButtons(T("Setup"), box, window)
	formPopup.Resize(fyne.NewSize(WindowWidth()*0.5, WindowHeight()*0.5))
	forms.FormToFyneForm(
		cubeSetupForm,
		box,
		formPopup,
		window,
		formSubmitCallback,
	)
	formPopup.Show()
}
