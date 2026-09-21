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

func TranslateContent(profile *models.CVProfile, source, target string, translateFn TranslateFn) (models.LocalizedContent, error) {
	profile.Normalize()

	sourceContent := profile.ContentForLang(source)
	if sourceContent.IsEmpty() {
		return models.LocalizedContent{}, nil
	}

	return translateContentWith(sourceContent, source, target, translateFn)
}

func ApplyTranslatedContent(profile *models.CVProfile, target string, translated models.LocalizedContent) {
	target = NormalizeLangCode(target)
	profile.SetContentForLang(target, translated)
	profile.Language = models.NormalizeLangCode(target)
	profile.Normalize()
}

func TranslateProfileWith(profile *models.CVProfile, source, target string, translateFn TranslateFn) error {
	translated, err := TranslateContent(profile, source, target, translateFn)
	if err != nil {
		return err
	}

	sourceContent := profile.ContentForLang(source)
	if sourceContent.IsEmpty() {
		if StorageSupported(target) {
			profile.Language = NormalizeLangCode(target)
			profile.NormalizeLanguage()
		}
		return nil
	}

	if StorageSupported(target) {
		ApplyTranslatedContent(profile, target, translated)
	}
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

	result := cloneLocalizedContent(content)

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
			j := j
			add(exp.Highlights[j], func(v string) { exp.Highlights[j] = v })
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

func cloneLocalizedContent(content *models.LocalizedContent) models.LocalizedContent {
	if content == nil {
		return models.EmptyLocalizedContent()
	}

	cloned := models.LocalizedContent{
		Summary:           content.Summary,
		PersonalTitle:     content.PersonalTitle,
		PersonalLanguages: content.PersonalLanguages,
		Experiences:       make([]models.Experience, len(content.Experiences)),
		Education:         append([]models.Education(nil), content.Education...),
		Projects:          append([]models.Project(nil), content.Projects...),
		SkillGroups:       make([]models.SkillGroup, len(content.SkillGroups)),
	}

	for i, exp := range content.Experiences {
		cloned.Experiences[i] = exp
		cloned.Experiences[i].Highlights = append([]string(nil), exp.Highlights...)
	}
	for i, group := range content.SkillGroups {
		cloned.SkillGroups[i] = group
		cloned.SkillGroups[i].Skills = append([]string(nil), group.Skills...)
	}

	return cloned
}
