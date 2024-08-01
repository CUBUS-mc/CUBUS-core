package translation

import (
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"path/filepath"
	"runtime"
)

var bundle *i18n.Bundle
var localizer *i18n.Localizer

// TODO: update components using T() when language is changed
func init() {
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}

	dir := filepath.Dir(filename)

	enPath := filepath.Join(dir, "locals", "active.en.toml")
	dePath := filepath.Join(dir, "locals", "active.de.toml")
	owoPath := filepath.Join(dir, "locals", "active.ky.toml")

	bundle.MustLoadMessageFile(enPath)
	bundle.MustLoadMessageFile(dePath)
	bundle.MustLoadMessageFile(owoPath)

	localizer = i18n.NewLocalizer(bundle, "en-US")
}

func T(key string) string {
	return localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: key,
		},
	})
}

func ChangeLanguage(lang string) {
	println("Changing language to", lang)
	if lang == "owo-UwU" {
		lang = "ky-KG"
	}
	localizer = i18n.NewLocalizer(bundle, lang)
}
