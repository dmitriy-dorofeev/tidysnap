package i18n

import (
	"os"
	"testing"
)

func TestSetLang(t *testing.T) {
	original := currentLang
	defer func() { currentLang = original }()

	SetLang("ru")
	if currentLang != "ru" {
		t.Errorf("expected ru, got %s", currentLang)
	}

	SetLang("en")
	if currentLang != "en" {
		t.Errorf("expected en, got %s", currentLang)
	}

	SetLang("RU")
	if currentLang != "ru" {
		t.Errorf("expected ru, got %s", currentLang)
	}

	SetLang("")
	// falls back to auto-detect; just ensure no panic
}

func TestT(t *testing.T) {
	original := currentLang
	defer func() { currentLang = original }()

	currentLang = "ru"
	if got := T("welcome_title"); got == "welcome_title" {
		t.Error("expected Russian translation for welcome_title")
	}

	currentLang = "en"
	if got := T("welcome_title"); got == "welcome_title" {
		t.Error("expected English translation for welcome_title")
	}

	currentLang = "unknown"
	if got := T("welcome_title"); got == "welcome_title" {
		t.Error("expected English fallback for unknown language")
	}
}

func TestNormalizeLang(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"ru_RU.UTF-8", "ru"},
		{"Russian", "ru"},
		{"en_US.UTF-8", "en"},
		{"ENG", "en"},
		{"fr_FR", "en"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := normalizeLang(tt.input); got != tt.want {
				t.Errorf("normalizeLang(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestDetectLangFromEnv(t *testing.T) {
	vars := []string{"LC_ALL", "LC_MESSAGES", "LC_CTYPE", "LANG"}
	orig := make(map[string]string)
	for _, v := range vars {
		orig[v] = os.Getenv(v)
		os.Unsetenv(v)
	}
	defer func() {
		for k, v := range orig {
			if v == "" {
				os.Unsetenv(k)
			} else {
				os.Setenv(k, v)
			}
		}
	}()

	os.Setenv("LANG", "ru_RU.UTF-8")
	if got := detectLang(); got != "ru" {
		t.Errorf("detectLang() = %q, want ru", got)
	}

	os.Setenv("LANG", "en_US.UTF-8")
	if got := detectLang(); got != "en" {
		t.Errorf("detectLang() = %q, want en", got)
	}
}
