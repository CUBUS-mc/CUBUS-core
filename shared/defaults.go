package shared

type UiType string

const (
	TUI  UiType = "tui"
	GUI  UiType = "gui"
	API  UiType = "api"
	NONE UiType = "none"
)

type Defaults struct {
	UI           UiType
	IconURL      string
	CubeAssetURL string
	Language     string
}

func NewDefaults() *Defaults {
	return &Defaults{
		UI:           GUI,
		IconURL:      "https://raw.githubusercontent.com/CUBUS-mc/CUBUS-core/master/assets/android.png",
		CubeAssetURL: "https://raw.githubusercontent.com/CUBUS-mc/CUBUS-core/master/assets/cube.svg",
		Language:     "en-US",
	}
}
