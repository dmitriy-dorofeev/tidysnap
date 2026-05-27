package i18n

import (
	"os"
	"strings"
)

var currentLang = detectLang()

func detectLang() string {
	for _, env := range []string{"LC_ALL", "LC_MESSAGES", "LANG"} {
		if v := os.Getenv(env); v != "" {
			v = strings.ToLower(v)
			if strings.HasPrefix(v, "ru") {
				return "ru"
			}
			return "en"
		}
	}
	return "en"
}

// T returns the translated string for the given key.
func T(key string) string {
	if m, ok := translations[currentLang]; ok {
		if s, ok := m[key]; ok {
			return s
		}
	}
	// Fallback to English.
	if m, ok := translations["en"]; ok {
		if s, ok := m[key]; ok {
			return s
		}
	}
	return key
}
