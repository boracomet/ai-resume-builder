package translate

import (
	"strings"

	"github.com/boradev/bora-cv/internal/models"
)

func TranslateProfile(profile *models.CVProfile, source, target, apiKey string) error {
	profile.Normalize()

	sourceContent := profile.ContentForLang(source)
	if sourceContent.IsEmpty() {
		profile.Language = target
		profile.NormalizeLanguage()
		return nil
	}

	translated, err := translateContent(sourceContent, source, target, apiKey)
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
		add(edu.School, func(v string) { edu.School = v })
		add(edu.StartDate, func(v string) { edu.StartDate = v })
		add(edu.EndDate, func(v string) { edu.EndDate = v })
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

	translated, err := TranslateTexts(texts, source, target, apiKey)
	if err != nil {
		return models.LocalizedContent{}, err
	}

	for i, value := range translated {
		apply[i](value)
	}

	return result, nil
}
