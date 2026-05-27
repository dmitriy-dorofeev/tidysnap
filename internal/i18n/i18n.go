package i18n

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var currentLang = detectLang()

func detectLang() string {
	for _, env := range []string{"LC_ALL", "LC_MESSAGES", "LC_CTYPE", "LANG"} {
		if v := os.Getenv(env); v != "" {
			return normalizeLang(v)
		}
	}
	if runtime.GOOS == "darwin" {
		if lang := macOSLang(); lang != "" {
			return lang
		}
	}
	return "en"
}

func normalizeLang(v string) string {
	v = strings.ToLower(v)
	if strings.HasPrefix(v, "ru") {
		return "ru"
	}
	return "en"
}

func macOSLang() string {
	out, err := exec.Command("defaults", "read", "-g", "AppleLanguages").Output()
	if err == nil {
		s := string(out)
		// Try to find first quoted language code
		start := strings.Index(s, `"`)
		if start != -1 {
			end := strings.Index(s[start+1:], `"`)
			if end != -1 {
				return normalizeLang(s[start+1 : start+1+end])
			}
		}
		// Fallback: try without quotes (e.g. (en, ru))
		start = strings.IndexFunc(s, func(r rune) bool {
			return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
		})
		if start != -1 {
			end := start + 1
			for end < len(s) && ((s[end] >= 'a' && s[end] <= 'z') || (s[end] >= 'A' && s[end] <= 'Z')) {
				end++
			}
			return normalizeLang(s[start:end])
		}
	}

	out, err = exec.Command("defaults", "read", "-g", "AppleLocale").Output()
	if err == nil {
		return normalizeLang(strings.TrimSpace(string(out)))
	}
	return ""
}

// SetLang explicitly sets the current language. Empty string falls back to auto-detection.
func SetLang(lang string) {
	switch strings.ToLower(lang) {
	case "ru", "russian":
		currentLang = "ru"
	case "en", "english":
		currentLang = "en"
	case "":
		currentLang = detectLang()
	default:
		currentLang = "en"
	}
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
