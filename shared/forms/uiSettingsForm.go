package forms

import "CUBUS-core/shared/translation"

func GetUiSettingsForm(
	currentLanguage string,
) *Form {
	form := NewForm(
		NewMultipleChoiceField(
			"language",
			[]DisplayCondition{&AlwaysDisplay{}},
			[]Validator{&ChoiceValidator{}},
			translation.T("Language"),
			translation.T("Choose your preferred language."),
			map[string]Option{
				"en-US": {
					Label:       translation.T("English"),
					Description: "",
				},
				"de-DE": {
					Label:       translation.T("Deutsch"),
					Description: "",
				},
				"owo-UwU": {
					Label:       translation.T("Engwish OwO"),
					Description: translation.T("UwUifyed English OwO"),
				},
			},
			currentLanguage,
		),
	)
	return form
}
