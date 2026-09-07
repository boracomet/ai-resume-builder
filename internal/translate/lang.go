package translate

import (
	"strings"
	"unicode"
)

// ValidLangCode reports whether code is a lowercase ISO 639-1 language code.
func ValidLangCode(code string) bool {
	code = strings.ToLower(strings.TrimSpace(code))
	if len(code) != 2 {
		return false
	}
	for _, r := range code {
		if r < 'a' || r > 'z' {
			return false
		}
	}
	return true
}

// NormalizeLangCode lowercases and trims a language code.
func NormalizeLangCode(code string) string {
	return strings.ToLower(strings.TrimSpace(code))
}

// StorageSupported reports whether translated CV content can be persisted in profile slots.
func StorageSupported(target string) bool {
	return ValidLangCode(NormalizeLangCode(target))
}

// ResolveSourceLang picks source language from explicit value or editing language fallback.
func ResolveSourceLang(explicit, editingLang, targetLang string) string {
	if ValidLangCode(explicit) {
		return NormalizeLangCode(explicit)
	}
	if ValidLangCode(editingLang) {
		return NormalizeLangCode(editingLang)
	}
	if targetLang == "tr" {
		return "en"
	}
	return "tr"
}

// MatchLangQuery returns true when a language matches a search query (code or names).
func MatchLangQuery(code, native, english, query string) bool {
	query = strings.TrimSpace(strings.ToLower(query))
	if query == "" {
		return true
	}
	if strings.Contains(code, query) {
		return true
	}
	if strings.Contains(strings.ToLower(native), query) {
		return true
	}
	if strings.Contains(strings.ToLower(english), query) {
		return true
	}
	return false
}

// FilterLanguages is used by tests; frontend has its own list.
func FilterLanguages(query string, codes []string, names map[string][2]string) []string {
	var matched []string
	for _, code := range codes {
		pair := names[code]
		if MatchLangQuery(code, pair[0], pair[1], query) {
			matched = append(matched, code)
		}
	}
	return matched
}

// IsLetter checks alphabetic runes for validation helpers.
func IsLetter(r rune) bool {
	return unicode.IsLetter(r)
}
