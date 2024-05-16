package translation

import (
	"testing"
)

func TestT(t *testing.T) {
	T := T
	// Set up a test case
	key := "Hello, world!"
	expected := "Hello, world!"

	// Call the function under test
	result := T(key)

	// Check the result
	if result != expected {
		t.Errorf("T(%q) = %q; want %q", key, result, expected)
	}

	// Test with a different language
	ChangeLanguage("de")
	expectedGerman := "Hallo, Welt!"
	result = T(key)

	if result != expectedGerman {
		t.Errorf("T(%q) = %q; want %q", key, result, expectedGerman)
	}
}
