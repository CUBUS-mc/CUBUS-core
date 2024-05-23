package forms

// TODO: Add translations

var CubeTypes = map[string]Option{
	"queen": {
		Label:       "Queen",
		Description: "The Queen is the central cube in the CUBUS network. It is the C2 server. There can only be one Queen in a CUBUS subnet.",
	},
	"security": {
		Label:       "Security",
		Description: "The Security cube is responsible for the security of the CUBUS network.",
	},
	"database": {
		Label:       "Database",
		Description: "The Database cube is responsible for storing data for the CUBUS network.",
	},
	"api": {
		Label:       "API",
		Description: "The API cube is responsible for providing an API for the CUBUS network.",
	},
	"cubus-mod": {
		Label:       "CUBUS Mod",
		Description: "This is a special cube that runs in a Minecraft mod.",
	},
	"discord-bot": {
		Label:       "Discord Bot",
		Description: "The Discord Bot cube is responsible for providing a Discord bot for the CUBUS network.",
	},
	"web": {
		Label:       "Web",
		Description: "The Web cube is responsible for providing a web interface for the CUBUS network.",
	},
	"drone": {
		Label:       "Drone",
		Description: "The Drone cubes are responsible for scanning, modifying and providing a proxy for the CUBUS network.",
	},
	"generic-worker": {
		Label:       "Generic Worker",
		Description: "The Generic Worker cube is a generic cube that can be used for any purpose by any other cube.",
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
			"Please select the location of the cube",
			"Cube Location",
			map[string]Option{
				"local": {
					Label:       "Local",
					Description: "Setup the cube on this device",
				},
				"remote": {
					Label:       "Remote",
					Description: "Setup the cube on a remote device",
				},
			},
			"",
		),
		NewTextField(
			"cubeName",
			[]DisplayCondition{&DisplayAfter{fieldId: "cubeLocation"}},
			[]Validator{&NotEmptyValidator{}},
			"Please enter the name of the cube",
			"Cube Name",
			"",
		),
		NewMultipleChoiceField(
			"cubeType",
			[]DisplayCondition{&DisplayAfter{fieldId: "cubeName"}},
			[]Validator{&ChoiceValidator{}},
			"Please select the type of the cube",
			"Cube Type",
			CubeTypes,
			"",
		),
		queenCubeSetupFields,
	)
	return form
}
