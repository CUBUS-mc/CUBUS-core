package tests

import (
	"CUBUS-core/translation"
	"testing"
)

func TestChangeLanguage(t *testing.T) {
	err := translation.ChangeLanguage("de")
	if err != nil {
		t.Errorf("ChangeLanguage() failed, expected nil, got %v", err)
	}
}

func TestLoadTranslations(t *testing.T) {
	err := translation.LoadTranslations()
	if err != nil {
		t.Errorf("LoadTranslations() failed, expected nil, got %v", err)
	}
}

func TestTranslate(t *testing.T) {
	_ = translation.LoadTranslations()
	result := translation.T("Test Translation")
	println(result)
	if result == "" {
		t.Errorf("Translate() failed, expected a string, got %v", result)
	}
}
