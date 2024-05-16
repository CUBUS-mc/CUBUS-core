package translation

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/message/catalog"
)

//go:generate gotext -srclang=en-US update -out=catalog.go -lang=en-US,de-DE CUBUS-core

// TODO: use "golang.org/x/text/language" "golang.org/x/text/message" "golang.org/x/text/message/catalog" to write the function t() that translates the english strings to the target language. If the input isn't english but a key, it should return the corresponding string in the target language.
// TODO: then add a method to change the target language
// TODO: The translations should be stored in a separate file and loaded at runtime
// TODO: replace all strings in the code with the t() function

var printer *message.Printer
var builder *catalog.Builder

func init() {
	builder = catalog.NewBuilder()
	err := LoadTranslations()
	if err != nil {
		return
	}
	printer = message.NewPrinter(language.English)
}

func T(key string) string {
	return printer.Sprintf(key)
}

func ChangeLanguage(lang string) error {
	tag, err := language.Parse(lang)
	if err != nil {
		return err
	}
	printer = message.NewPrinter(tag)
	return nil
}

func LoadTranslations() error {
	return nil
}
