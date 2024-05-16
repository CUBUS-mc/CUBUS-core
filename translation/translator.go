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

func init() {
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	// Get the current file's path
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}

	// Get the directory of the current file
	dir := filepath.Dir(filename)

	// Construct the absolute paths to the TOML files
	enPath := filepath.Join(dir, "locals", "active.en.toml")
	dePath := filepath.Join(dir, "locals", "active.de.toml")

	bundle.MustLoadMessageFile(enPath)
	bundle.MustLoadMessageFile(dePath)

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
	localizer = i18n.NewLocalizer(bundle, lang)
}
