package main

import (
	"encoding/json"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"os"
)

func loadConfig() map[string]interface{} {
	if !checkPathExists("config.json") {
		configureGUI()
	}

	configString := readFile("config.json")

	var config map[string]interface{}
	err := json.Unmarshal([]byte(configString), &config)
	if err != nil {
		panic(err)
	}

	return config
}

func reconfigure() {
	if _, err := os.Stat("config.json"); os.IsNotExist(err) {
		configureGUI()
		return
	}

	var config map[string]interface{}
	err := json.Unmarshal([]byte(readFile("config.json")), &config)
	if err != nil {
		panic(err)
	}

	switch config["ui_type"].(string) {
	case "c":
		configureCLI()
	case "g":
		configureGUI()
	default:
		configureGUI()
	}
}

func configureCLI() {
	print(banner)
	print("                              CONFIGURATOR\n\n")
	var config = make(map[string]interface{})

	walkTreeCLI(configDialogueTree, config, None)

	configString, err := json.Marshal(config)
	if err != nil {
		panic(err)
	}

	writeFile("config.json", string(configString))
}

func configureGUI() {
	configApp := app.New()
	configWindow := configApp.NewWindow("QUBUS Configurator")
	configWindow.Resize(fyne.NewSize(500, 320))
	configWindow.CenterOnScreen()
	configWindow.SetFixedSize(true)
	configWindow.SetCloseIntercept(func() {
		dialog.NewConfirm("Are you sure?", "Are you sure you want to exit the configurator? Progress will be lost.", func(b bool) {
			if b {
				os.Exit(0)
			}
		}, configWindow).Show()
	})
	theme.InnerPadding()

	heading := widget.NewLabel("QUBUS Configurator")
	heading.Alignment = fyne.TextAlignCenter
	heading.TextStyle.Bold = true

	config := make(map[string]interface{})

	okButton := widget.NewButton("OK", func() {
		configString, err := json.Marshal(config)
		if err != nil {
			panic(err)
		}
		writeFile("config.json", string(configString))
		configWindow.Close()
	})

	configContainer := container.NewVBox()
	walkTreeGUI(configDialogueTree, config, configContainer)
	configWindow.SetContent(container.NewPadded(container.NewVBox(heading, container.NewHScroll(configContainer), layout.NewSpacer(), okButton)))

	configWindow.ShowAndRun()
}
