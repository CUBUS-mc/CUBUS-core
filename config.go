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

// TODO: Load a configuration tree that shows what questions should be asked and what will follow from the answers (E.g. Cube type -> Drone -> Drone options or Cube type -> Queen -> Queen options)

func loadConfig() map[string]interface{} {
	if _, err := os.Stat("config.json"); os.IsNotExist(err) {
		configureCLI()
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
	configWindow.Resize(fyne.NewSize(420, 320))
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
	var pages []*fyne.Container

	cubeTypeLabel := widget.NewLabel("Cube type:")
	cubeTypes := []string{"Queen", "Security", "Drone", "Database", "API", "Website", "Discord bot"}
	cubeTypeDropdown := widget.NewSelect(cubeTypes, func(input string) { // TODO: Remove old page if reselecting; add page if default is left so no change
		var newPage *fyne.Container
		switch input {
		case "Queen":
			newPage = container.New(layout.NewFormLayout(), widget.NewLabel("Queen options"))
		case "Security":
			newPage = container.New(layout.NewFormLayout(), widget.NewLabel("Security options"))
		case "Drone":
			newPage = container.New(layout.NewFormLayout(), widget.NewLabel("Drone options"))
		case "Database":
			newPage = container.New(layout.NewFormLayout(), widget.NewLabel("Database options"))
		case "API":
			newPage = container.New(layout.NewFormLayout(), widget.NewLabel("API options"))
		case "Website":
			newPage = container.New(layout.NewFormLayout(), widget.NewLabel("Website options"))
		case "Discord bot":
			newPage = container.New(layout.NewFormLayout(), widget.NewLabel("Discord bot options"))
		}
		if len(pages) > 1 {
			pages = append(pages[:len(pages)-1], append([]*fyne.Container{newPage}, pages[len(pages)-1:]...)...)
		} else {
			pages = append(pages, newPage)
		}
	})
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

	pages = []*fyne.Container{container.New(layout.NewFormLayout(), cubeTypeLabel, cubeTypeDropdown), container.New(layout.NewFormLayout(), uiTypeLabel, uiTypeDropdown)}
	var activePage = pages[0]

	var nextButton *widget.Button
	nextButton = widget.NewButton("Next", func() { // TODO: Add back button
		page++
		if page >= len(pages)-1 {
			okButton.Show()
			nextButton.Hide()
		}
		activePage = pages[page]
		configWindow.SetContent(container.NewPadded(container.NewVBox(heading, layout.NewSpacer(), activePage, layout.NewSpacer(), nextButton, okButton)))
	})

	configWindow.SetContent(container.NewPadded(container.NewVBox(heading, layout.NewSpacer(), activePage, layout.NewSpacer(), nextButton, okButton)))

	configWindow.ShowAndRun()
}
