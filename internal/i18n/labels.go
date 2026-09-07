package i18n

import "strings"

var labels = map[string]map[string]string{
	"tr": {
		"summary":    "ÖZET",
		"experience": "DENEYİM",
		"education":  "EĞİTİM",
		"projects":   "PROJELER",
		"skills":     "BECERİLER",
		"birthDate":  "Doğum",
		"languages":  "Dil",
		"photoAlt":   "profil fotoğrafı",
	},
	"en": {
		"summary":    "SUMMARY",
		"experience": "EXPERIENCE",
		"education":  "EDUCATION",
		"projects":   "PROJECTS",
		"skills":     "SKILLS",
		"birthDate":  "Birth",
		"languages":  "Language",
		"photoAlt":   "profile photo",
	},
}

func NormalizeLang(lang string) string {
	lang = strings.ToLower(strings.TrimSpace(lang))
	if lang == "en" {
		return "en"
	}
	return "tr"
}

func T(lang, key string) string {
	lang = NormalizeLang(lang)
	if m, ok := labels[lang]; ok {
		if v, ok := m[key]; ok {
			return v
		}
	}
	return key
}
