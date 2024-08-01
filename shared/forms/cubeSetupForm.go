package forms

import (
	"CUBUS-core/shared/translation"
)

var CubeTypes = map[string]Option{
	"queen": {
		Label:       translation.T("Queen"),
		Description: translation.T("The Queen is the central cube in the CUBUS network. It is the C2 server. There can only be one Queen in a CUBUS subnet."),
	},
	"security": {
		Label:       translation.T("Security"),
		Description: translation.T("The Security cube is responsible for the security of the CUBUS network."),
	},
	"database": {
		Label:       translation.T("Database"),
		Description: translation.T("The Database cube is responsible for storing data for the CUBUS network."),
	},
	"api": {
		Label:       translation.T("API"),
		Description: translation.T("The API cube is responsible for providing an API for the CUBUS network."),
	},
	"cubus-mod": {
		Label:       translation.T("CUBUS Mod"),
		Description: translation.T("This is a special cube that runs in a Minecraft mod."),
	},
	"discord-bot": {
		Label:       translation.T("Discord Bot"),
		Description: translation.T("The Discord Bot cube is responsible for providing a Discord bot for the CUBUS network."),
	},
	"web": {
		Label:       translation.T("Web"),
		Description: translation.T("The Web cube is responsible for providing a web interface for the CUBUS network."),
	},
	"drone": {
		Label:       translation.T("Drone"),
		Description: translation.T("The Drone cubes are responsible for scanning, modifying and providing a proxy for the CUBUS network."),
	},
	"generic-worker": {
		Label:       translation.T("Generic Worker"),
		Description: translation.T("The Generic Worker cube is a generic cube that can be used for any purpose by any other cube."),
	},
}

func GetCubeSetupForm() *Form {
	queenCubeSetupFields := NewFieldGroup(
		"queenCubeSetup",
		[]DisplayCondition{&HasValueDisplayCondition{fieldId: "cubeType", value: "queen"}},
		[]Validator{&IsValidValidator{fieldIds: []string{"queenName"}}},
		"",
		NewMessage(
			"queenCubeSetupMessage",
			[]DisplayCondition{&AlwaysDisplay{}},
			"CUBUS Queen setup placeholder",
		),
	)

	form := NewForm(
		NewMultipleChoiceField(
			"cubeLocation",
			[]DisplayCondition{&AlwaysDisplay{}},
			[]Validator{&ChoiceValidator{}},
			translation.T("Please select the location of the cube"),
			translation.T("Cube Location"),
			map[string]Option{
				"local": {
					Label:       translation.T("Local"),
					Description: translation.T("Setup the cube on this device"),
				},
				"remote": {
					Label:       translation.T("Remote"),
					Description: translation.T("Setup the cube on a remote device"),
				},
			},
			"",
		),
		NewUrlField(
			"remoteUrl",
			[]DisplayCondition{&DisplayAfter{fieldId: "cubeLocation"}, &HasValueDisplayCondition{fieldId: "cubeLocation", value: "remote"}},
			[]Validator{&UrlValidator{}},
			translation.T("Please enter the URL of the remote CUBUS-Core orchestrator server"),
			translation.T("Remote URL"),
			"",
		),
		NewTextField(
			"cubeName",
			[]DisplayCondition{&OrDisplayCondition{[]DisplayCondition{&HasValueDisplayCondition{fieldId: "cubeLocation", value: "local"}, &DisplayAfter{fieldId: "remoteUrl"}}}},
			[]Validator{&NotEmptyValidator{}},
			translation.T("Please enter the name of the cube"),
			translation.T("Cube Name"),
			"",
		),
		NewMultipleChoiceField(
			"cubeType",
			[]DisplayCondition{&DisplayAfter{fieldId: "cubeName"}},
			[]Validator{&ChoiceValidator{}},
			translation.T("Please select the type of the cube"),
			translation.T("Cube Type"),
			CubeTypes,
			"",
		),
		queenCubeSetupFields,
	)
	return form
}
