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
	"strings"
)

func loadConfig() map[string]interface{} {
	if _, err := os.Stat("config.json"); os.IsNotExist(err) {
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
	var config map[string]interface{}
	config = make(map[string]interface{})

	println("What type of UI do you want to use?")
	println("c - CLI, g - GUI (default: GUI)")
	config["ui_type"] = strings.ToLower(askForString("Type of UI: ", "g", func(input string) bool {
		switch strings.ToLower(input) {
		case "c", "g":
			return true
		default:
			return false
		}
	}))

	println("What type of cube do you want to setup?")
	println("q - Queen, s - Security, d - Drone, db - Database, a - API, w - Website, b - Discord bot (default: Drone)")
	config["cube_type"] = strings.ToLower(askForString("Type of cube: ", "d", func(input string) bool {
		switch strings.ToLower(input) {
		case "q", "s", "d", "db", "a", "w", "b":
			return true
		default:
			return false
		}
	}))

	configString, err := json.Marshal(config)
	if err != nil {
		panic(err)
	}

	writeFile("config.json", string(configString))
	print(banner)
}

func configureGUI() {
	configApp := app.New()
	configWindow := configApp.NewWindow("QUBUS Configurator")
	configWindow.Resize(fyne.NewSize(400, 300))
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

	var page = 0

	cubeTypeLabel := widget.NewLabel("Cube type:")
	cubeTypes := []string{"Queen", "Security", "Drone", "Database", "API", "Website", "Discord bot"}
	cubeTypeDropdown := widget.NewSelect(cubeTypes, func(string) {})
	cubeTypeDropdown.SetSelected("Drone")

	uiTypeLabel := widget.NewLabel("UI type:")
	uiTypes := []string{"CLI", "GUI"}
	uiTypeDropdown := widget.NewSelect(uiTypes, func(string) {})
	uiTypeDropdown.SetSelected("GUI")

	okButton := widget.NewButton("OK", func() {
		config := make(map[string]interface{})
		config["ui_type"] = map[string]string{"CLI": "c", "GUI": "g"}[uiTypeDropdown.Selected]
		config["cube_type"] = map[string]string{"Queen": "q", "Security": "s", "Drone": "d", "Database": "db", "API": "a", "Website": "w", "Discord bot": "b"}[cubeTypeDropdown.Selected]
		configString, err := json.Marshal(config)
		if err != nil {
			panic(err)
		}
		writeFile("config.json", string(configString))
		configWindow.Close()
	})
	okButton.Alignment = widget.ButtonAlignCenter
	okButton.Hide()
	okButton.Move(fyne.NewPos(0, 300-32))

	var pages = []*fyne.Container{container.New(layout.NewFormLayout(), cubeTypeLabel, cubeTypeDropdown), container.New(layout.NewFormLayout(), uiTypeLabel, uiTypeDropdown)}
	var activePage = pages[0]

	// TODO: Add Next button to go to the next page of the configurator with more options so that options can be shown/hidden based on the cube type
	var nextButton *widget.Button
	nextButton = widget.NewButton("Next", func() {
		page++ // TODO: add animation
		if page >= len(pages)-1 {
			okButton.Show()
			nextButton.Hide()
			return
		}
		activePage = pages[page]
		configWindow.SetContent(container.NewPadded(container.NewVBox(heading, activePage, nextButton, okButton)))
	})

	configWindow.SetContent(container.NewPadded(container.NewVBox(heading, activePage, nextButton, okButton)))

	configWindow.ShowAndRun()
}
