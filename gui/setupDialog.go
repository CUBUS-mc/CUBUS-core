package gui

import (
	"CUBUS-core/orchestrator/client"
	"CUBUS-core/shared"
	"CUBUS-core/shared/forms"
	"CUBUS-core/shared/translation"
	"CUBUS-core/shared/types"
	"CUBUS-core/shared/types/gui"
	"context"
	"crypto/rsa"
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
	addServerUrlIfNotExists func(serverUrl string),
) {
	T := translation.T

	cubeConfig := types.CubeConfig{
		Id:          uuid.New().String(),
		CubeName:    "",
		PublicKey:   rsa.PublicKey{},
		QueueServer: types.QueueServerConfig{},
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
			}
		}

		var serverUrl string
		if values["cubeLocation"] == "local" {
			serverUrl = "localhost:25560"
		} else {
			serverUrl = values["remoteUrl"]
			addServerUrlIfNotExists(serverUrl)
		}
		orchestratorClient, err := client.NewClient(serverUrl)
		if err != nil {
			dialog.ShowError(err, window)
			return
		}
		// TODO: THIS IS A TEMPORARY FIX
		cubeConfig.QueueServer = types.QueueServerConfig{
			Url:      "localhost:6379",
			Username: "",
			Password: "",
			DB:       0,
		}
		_, err = orchestratorClient.CreateCube(context.Background(), cubeConfig)
		if err != nil {
			dialog.ShowError(err, window)
			return
		}
		println("Created cube with id: ", cubeConfig.Id)
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
