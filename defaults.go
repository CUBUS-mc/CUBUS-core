package main

type UiType string

const (
	CLI UiType = "cli"
	GUI UiType = "gui"
	API UiType = "api"
)

type Defaults struct {
	UI           UiType
	IconURL      string
	CubeAssetURL string
}

func NewDefaults() *Defaults {
	return &Defaults{
		UI:           GUI,
		IconURL:      "https://raw.githubusercontent.com/CUBUS-mc/CUBUS-core/master/assets/android.png",
		CubeAssetURL: "https://raw.githubusercontent.com/CUBUS-mc/CUBUS-core/master/assets/cube.svg",
	}
}
