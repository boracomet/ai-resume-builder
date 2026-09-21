package translate

import (
	"testing"

	"github.com/boracomet/ai-resume-builder/internal/models"
)

func TestTranslateProfileDoesNotMutateSourceContent(t *testing.T) {
	profile := &models.CVProfile{
		Language: "tr",
		Personal: models.PersonalInfo{
			Name:     "Bora Ata Türkoğlu",
			Location: "İstanbul",
		},
		ContentTR: models.LocalizedContent{
			Summary:           "Türkçe özet",
			PersonalTitle:     "Yazılım Geliştirici",
			PersonalLanguages: "Türkçe (Ana Dil), İngilizce (B2)",
			Experiences: []models.Experience{
				{
					Title:       "Web Geliştirici",
					Company:     "Örnek Şirket",
					StartDate:   "Ocak 2023",
					EndDate:     "Mart 2025",
					Duration:    "26 Ay",
					Location:    "Uzaktan",
					Description: "Şirket açıklaması",
					Highlights:  []string{"Ekip koordinasyonuna destek oldum.", "Dağıtım süreçlerine katkı sağladım."},
				},
			},
			Education: []models.Education{
				{
					Degree:      "Web Tasarımı ve Kodlama",
					Institution: "Anadolu Üniversitesi",
					StartDate:   "2022",
					EndDate:     "Devam Ediyor",
					Description: "Ön lisans eğitimi",
				},
			},
			Projects: []models.Project{
				{Name: "Kişisel Site", URL: "https://example.com", Description: "Portfolyo sitesi"},
			},
			SkillGroups: []models.SkillGroup{
				{Category: "Diller", Skills: []string{"Go", "JavaScript"}},
			},
		},
	}

	err := TranslateProfileWith(profile, "tr", "en", func(texts []string, source, target string) ([]string, error) {
		out := make([]string, len(texts))
		for i, text := range texts {
			out[i] = "EN:" + text
		}
		return out, nil
	})
	if err != nil {
		t.Fatalf("TranslateProfileWith() error = %v", err)
	}

	tr := profile.ContentTR
	en := profile.ContentEN

	assertUnchanged := func(name, got, want string) {
		t.Helper()
		if got != want {
			t.Fatalf("source %s mutated: got %q want %q", name, got, want)
		}
	}
	assertTranslated := func(name, got, original string) {
		t.Helper()
		if got != "EN:"+original {
			t.Fatalf("translated %s = %q, want %q", name, got, "EN:"+original)
		}
	}

	assertUnchanged("summary", tr.Summary, "Türkçe özet")
	assertUnchanged("personalTitle", tr.PersonalTitle, "Yazılım Geliştirici")
	assertUnchanged("personalLanguages", tr.PersonalLanguages, "Türkçe (Ana Dil), İngilizce (B2)")
	assertUnchanged("exp.title", tr.Experiences[0].Title, "Web Geliştirici")
	assertUnchanged("exp.company", tr.Experiences[0].Company, "Örnek Şirket")
	assertUnchanged("exp.startDate", tr.Experiences[0].StartDate, "Ocak 2023")
	assertUnchanged("exp.endDate", tr.Experiences[0].EndDate, "Mart 2025")
	assertUnchanged("exp.duration", tr.Experiences[0].Duration, "26 Ay")
	assertUnchanged("exp.location", tr.Experiences[0].Location, "Uzaktan")
	assertUnchanged("exp.description", tr.Experiences[0].Description, "Şirket açıklaması")
	assertUnchanged("exp.highlights[0]", tr.Experiences[0].Highlights[0], "Ekip koordinasyonuna destek oldum.")
	assertUnchanged("exp.highlights[1]", tr.Experiences[0].Highlights[1], "Dağıtım süreçlerine katkı sağladım.")
	assertUnchanged("edu.degree", tr.Education[0].Degree, "Web Tasarımı ve Kodlama")
	assertUnchanged("edu.institution", tr.Education[0].Institution, "Anadolu Üniversitesi")
	assertUnchanged("edu.startDate", tr.Education[0].StartDate, "2022")
	assertUnchanged("edu.endDate", tr.Education[0].EndDate, "Devam Ediyor")
	assertUnchanged("edu.description", tr.Education[0].Description, "Ön lisans eğitimi")
	assertUnchanged("project.name", tr.Projects[0].Name, "Kişisel Site")
	assertUnchanged("project.description", tr.Projects[0].Description, "Portfolyo sitesi")
	assertUnchanged("skill.category", tr.SkillGroups[0].Category, "Diller")
	assertUnchanged("skill.skills[0]", tr.SkillGroups[0].Skills[0], "Go")
	assertUnchanged("personal.name", profile.Personal.Name, "Bora Ata Türkoğlu")
	assertUnchanged("personal.location", profile.Personal.Location, "İstanbul")

	assertTranslated("summary", en.Summary, "Türkçe özet")
	assertTranslated("personalTitle", en.PersonalTitle, "Yazılım Geliştirici")
	assertTranslated("personalLanguages", en.PersonalLanguages, "Türkçe (Ana Dil), İngilizce (B2)")
	assertTranslated("exp.title", en.Experiences[0].Title, "Web Geliştirici")
	assertTranslated("exp.company", en.Experiences[0].Company, "Örnek Şirket")
	assertTranslated("exp.highlights[0]", en.Experiences[0].Highlights[0], "Ekip koordinasyonuna destek oldum.")
	assertTranslated("edu.degree", en.Education[0].Degree, "Web Tasarımı ve Kodlama")
	if en.Education[0].EndDate != "Ongoing" {
		t.Fatalf("edu.endDate mapped = %q, want Ongoing", en.Education[0].EndDate)
	}
	assertTranslated("project.description", en.Projects[0].Description, "Portfolyo sitesi")
	assertTranslated("skill.category", en.SkillGroups[0].Category, "Diller")
	if en.SkillGroups[0].Skills[0] != "Go" {
		t.Fatalf("skill items should pass through, got %q", en.SkillGroups[0].Skills[0])
	}

	en.Experiences[0].Highlights[0] = "MUTATED-EN"
	en.SkillGroups[0].Skills[0] = "MUTATED-SKILL"
	if tr.Experiences[0].Highlights[0] == "MUTATED-EN" {
		t.Fatal("source highlights still share backing array with translated content")
	}
	if tr.SkillGroups[0].Skills[0] == "MUTATED-SKILL" {
		t.Fatal("source skills still share backing array with translated content")
	}
}
