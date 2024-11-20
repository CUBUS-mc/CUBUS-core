package forms

import (
	"CUBUS-core/shared/translation"
)

func GetCubeSetupForm() *Form {
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
		// TODO: Add mutiple choice field for selecting the queue server
	)
	return form
}
