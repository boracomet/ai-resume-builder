package translate

import (
	"regexp"
	"strings"

	"github.com/boracomet/ai-resume-builder/internal/models"
	"github.com/boracomet/ai-resume-builder/internal/openaiusage"
)

var yearOnlyRe = regexp.MustCompile(`^\d{4}$`)

var educationStatusTranslations = map[string]map[string]string{
	"devam ediyor": {"en": "Ongoing"},
	"ongoing":      {"tr": "Devam Ediyor"},
	"present":      {"tr": "Devam Ediyor"},
}

type TranslateFn func(texts []string, source, target string) ([]string, error)

func TranslateProfile(profile *models.CVProfile, source, target, apiKey string) error {
	return TranslateProfileWith(profile, source, target, func(texts []string, src, tgt string) ([]string, error) {
		return TranslateTexts(texts, src, tgt, apiKey)
	})
}

func TranslateProfileWithOpenAI(profile *models.CVProfile, source, target, apiKey, model string) (openaiusage.Usage, float64, error) {
	var totalUsage openaiusage.Usage
	var totalCost float64
	err := TranslateProfileWith(profile, source, target, func(texts []string, src, tgt string) ([]string, error) {
		translated, usage, cost, err := TranslateTextsWithOpenAI(texts, src, tgt, apiKey, model)
		if err != nil {
			return nil, err
		}
		totalUsage = totalUsage.Add(usage)
		totalCost += cost
		return translated, nil
	})
	return totalUsage, totalCost, err
}

func TranslateProfileWith(profile *models.CVProfile, source, target string, translateFn TranslateFn) error {
	profile.Normalize()

	sourceContent := profile.ContentForLang(source)
	if sourceContent.IsEmpty() {
		profile.Language = target
		profile.NormalizeLanguage()
		return nil
	}

	translated, err := translateContentWith(sourceContent, source, target, translateFn)
	if err != nil {
		return err
	}

	if target == "en" {
		profile.ContentEN = translated
	} else {
		profile.ContentTR = translated
	}

	profile.Language = target
	profile.NormalizeLanguage()
	return nil
}

func translateContent(content *models.LocalizedContent, source, target, apiKey string) (models.LocalizedContent, error) {
	return translateContentWith(content, source, target, func(texts []string, src, tgt string) ([]string, error) {
		return TranslateTexts(texts, src, tgt, apiKey)
	})
}

func translateContentWith(content *models.LocalizedContent, source, target string, translateFn TranslateFn) (models.LocalizedContent, error) {
	var texts []string
	var apply []func(string)

	add := func(value string, setter func(string)) {
		value = strings.TrimSpace(value)
		if value == "" {
			return
		}
		texts = append(texts, value)
		apply = append(apply, setter)
	}

	result := models.LocalizedContent{
		Experiences: append([]models.Experience(nil), content.Experiences...),
		Education:   append([]models.Education(nil), content.Education...),
		Projects:    append([]models.Project(nil), content.Projects...),
		SkillGroups: append([]models.SkillGroup(nil), content.SkillGroups...),
	}

	add(content.Summary, func(v string) { result.Summary = v })
	add(content.PersonalTitle, func(v string) { result.PersonalTitle = v })
	add(content.PersonalLanguages, func(v string) { result.PersonalLanguages = v })

	for i := range result.Experiences {
		exp := &result.Experiences[i]
		add(exp.Title, func(v string) { exp.Title = v })
		add(exp.Company, func(v string) { exp.Company = v })
		add(exp.StartDate, func(v string) { exp.StartDate = v })
		add(exp.EndDate, func(v string) { exp.EndDate = v })
		add(exp.Duration, func(v string) { exp.Duration = v })
		add(exp.Location, func(v string) { exp.Location = v })
		add(exp.Description, func(v string) { exp.Description = v })
		for j := range exp.Highlights {
			highlight := &exp.Highlights[j]
			add(*highlight, func(v string) { *highlight = v })
		}
	}

	for i := range result.Education {
		edu := &result.Education[i]
		add(edu.Degree, func(v string) { edu.Degree = v })
		add(edu.Institution, func(v string) { edu.Institution = v })
		if mapped, ok := translateEducationScalar(edu.StartDate, target); ok {
			edu.StartDate = mapped
		} else {
			add(edu.StartDate, func(v string) { edu.StartDate = v })
		}
		if mapped, ok := translateEducationScalar(edu.EndDate, target); ok {
			edu.EndDate = mapped
		} else {
			add(edu.EndDate, func(v string) { edu.EndDate = v })
		}
		add(edu.Description, func(v string) { edu.Description = v })
	}

	for i := range result.Projects {
		project := &result.Projects[i]
		add(project.Name, func(v string) { project.Name = v })
		add(project.Description, func(v string) { project.Description = v })
	}

	for i := range result.SkillGroups {
		group := &result.SkillGroups[i]
		add(group.Category, func(v string) { group.Category = v })
	}

	if len(texts) == 0 {
		return result, nil
	}

	translated, err := translateFn(texts, source, target)
	if err != nil {
		return models.LocalizedContent{}, err
	}

	for i, value := range translated {
		apply[i](value)
	}

	return result, nil
}

func translateEducationScalar(value, target string) (string, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", false
	}
	if yearOnlyRe.MatchString(value) {
		return value, true
	}
	key := strings.ToLower(value)
	if translations, ok := educationStatusTranslations[key]; ok {
		if mapped, ok := translations[target]; ok {
			return mapped, true
		}
	}
	return "", false
}
