package models

import (
	"encoding/json"
	"strings"
	"time"
)

type PersonalInfo struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	LinkedIn  string `json:"linkedin"`
	GitHub    string `json:"github"`
	Portfolio string `json:"portfolio"`
	Location  string `json:"location"`
	BirthDate string `json:"birthDate"`
	// Legacy fields kept for migration from old single-language profiles.
	Title     string `json:"title,omitempty"`
	Languages string `json:"languages,omitempty"`
}

type Experience struct {
	Title       string   `json:"title"`
	Company     string   `json:"company"`
	StartDate   string   `json:"startDate"`
	EndDate     string   `json:"endDate"`
	Duration    string   `json:"duration"`
	Location    string   `json:"location"`
	Description string   `json:"description"`
	Highlights  []string `json:"highlights"`
}

type Education struct {
	Degree      string `json:"degree"`
	Institution string `json:"institution"`
	StartDate   string `json:"startDate"`
	EndDate     string `json:"endDate"`
	Description string `json:"description"`
}

func (e *Education) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	unmarshalField := func(key string, dest *string) {
		if value, ok := raw[key]; ok {
			_ = json.Unmarshal(value, dest)
		}
	}
	unmarshalField("degree", &e.Degree)
	unmarshalField("institution", &e.Institution)
	unmarshalField("startDate", &e.StartDate)
	unmarshalField("endDate", &e.EndDate)
	unmarshalField("description", &e.Description)
	if strings.TrimSpace(e.Institution) == "" {
		unmarshalField("school", &e.Institution)
	}
	return nil
}

func (e Education) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		"degree":      e.Degree,
		"institution": e.Institution,
		"startDate":   e.StartDate,
		"endDate":     e.EndDate,
		"description": e.Description,
	})
}

type Project struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	Description string `json:"description"`
}

type SkillGroup struct {
	Category string   `json:"category"`
	Skills   []string `json:"skills"`
}

type LocalizedContent struct {
	Summary           string       `json:"summary"`
	Experiences       []Experience `json:"experiences"`
	Education         []Education  `json:"education"`
	Projects          []Project    `json:"projects"`
	SkillGroups       []SkillGroup `json:"skillGroups"`
	PersonalTitle     string       `json:"personalTitle"`
	PersonalLanguages string       `json:"personalLanguages"`
}

const (
	DefaultPhotoSize        = 96
	DefaultPhotoBorderWidth = 2
	DefaultPhotoBorderColor = "#2563eb"
)

type CVProfile struct {
	ID               int64                        `json:"id"`
	Name             string                       `json:"name"`
	Language         string                       `json:"language"`
	Personal         PersonalInfo                 `json:"personal"`
	ContentTR        LocalizedContent             `json:"contentTR"`
	ContentEN        LocalizedContent             `json:"contentEN"`
	ContentLanguages map[string]*LocalizedContent `json:"contentLanguages,omitempty"`
	PhotoBase64      string           `json:"photoBase64"`
	PhotoSize        int              `json:"photoSize"`
	PhotoBorderWidth int              `json:"photoBorderWidth"`
	PhotoBorderColor string           `json:"photoBorderColor"`
	CreatedAt        time.Time        `json:"createdAt"`
	UpdatedAt        time.Time        `json:"updatedAt"`

	// Legacy single-language fields (migrated into ContentTR on load).
	Summary     string       `json:"summary,omitempty"`
	Experiences []Experience `json:"experiences,omitempty"`
	Education   []Education  `json:"education,omitempty"`
	Projects    []Project    `json:"projects,omitempty"`
	SkillGroups []SkillGroup `json:"skillGroups,omitempty"`
}

func EmptyLocalizedContent() LocalizedContent {
	return LocalizedContent{
		Experiences: []Experience{},
		Education:   []Education{},
		Projects:    []Project{},
		SkillGroups: []SkillGroup{},
	}
}

func (c *LocalizedContent) IsEmpty() bool {
	if c == nil {
		return true
	}
	if strings.TrimSpace(c.Summary) != "" ||
		strings.TrimSpace(c.PersonalTitle) != "" ||
		strings.TrimSpace(c.PersonalLanguages) != "" {
		return false
	}
	return len(c.Experiences) == 0 &&
		len(c.Education) == 0 &&
		len(c.Projects) == 0 &&
		len(c.SkillGroups) == 0
}

func ValidLangCode(lang string) bool {
	lang = strings.ToLower(strings.TrimSpace(lang))
	if len(lang) != 2 {
		return false
	}
	for _, r := range lang {
		if r < 'a' || r > 'z' {
			return false
		}
	}
	return true
}

func (p *CVProfile) ContentForLang(lang string) *LocalizedContent {
	lang = NormalizeLangCode(lang)
	if lang == "en" {
		return &p.ContentEN
	}
	if lang == "tr" {
		return &p.ContentTR
	}
	if p.ContentLanguages == nil {
		p.ContentLanguages = map[string]*LocalizedContent{}
	}
	if p.ContentLanguages[lang] == nil {
		p.ContentLanguages[lang] = ptrLocalizedContent(EmptyLocalizedContent())
	}
	return p.ContentLanguages[lang]
}

func ptrLocalizedContent(content LocalizedContent) *LocalizedContent {
	return &content
}

func (p *CVProfile) SetContentForLang(lang string, content LocalizedContent) {
	lang = NormalizeLangCode(lang)
	if lang == "en" {
		p.ContentEN = content
		return
	}
	if lang == "tr" {
		p.ContentTR = content
		return
	}
	if p.ContentLanguages == nil {
		p.ContentLanguages = map[string]*LocalizedContent{}
	}
	p.ContentLanguages[lang] = ptrLocalizedContent(content)
}

func (p *CVProfile) ActiveContent() *LocalizedContent {
	return p.ContentForLang(p.Language)
}

func (p *CVProfile) DisplayTitle() string {
	return strings.TrimSpace(p.ActiveContent().PersonalTitle)
}

func (p *CVProfile) DisplayLanguages() string {
	return strings.TrimSpace(p.ActiveContent().PersonalLanguages)
}

func (p *CVProfile) DisplaySummary() string {
	return p.ActiveContent().Summary
}

func (p *CVProfile) DisplayExperiences() []Experience {
	content := p.ActiveContent()
	if content.Experiences == nil {
		return []Experience{}
	}
	return content.Experiences
}

func (p *CVProfile) DisplayEducation() []Education {
	content := p.ActiveContent()
	if content.Education == nil {
		return []Education{}
	}
	return content.Education
}

func (p *CVProfile) DisplayProjects() []Project {
	content := p.ActiveContent()
	if content.Projects == nil {
		return []Project{}
	}
	return content.Projects
}

func (p *CVProfile) DisplaySkillGroups() []SkillGroup {
	content := p.ActiveContent()
	if content.SkillGroups == nil {
		return []SkillGroup{}
	}
	return content.SkillGroups
}

func NormalizeLangCode(lang string) string {
	lang = strings.ToLower(strings.TrimSpace(lang))
	if lang == "en" {
		return "en"
	}
	if ValidLangCode(lang) {
		return lang
	}
	return "tr"
}

func (p *CVProfile) NormalizeLanguage() {
	p.Language = NormalizeLangCode(p.Language)
}

func (p *CVProfile) Normalize() {
	p.MigrateFromLegacy()
	p.NormalizePhotoSettings()
	p.NormalizeLanguage()
	p.ensureContentSlices()
}

func (p *CVProfile) ensureContentSlices() {
	if p.ContentTR.Experiences == nil {
		p.ContentTR.Experiences = []Experience{}
	}
	if p.ContentTR.Education == nil {
		p.ContentTR.Education = []Education{}
	}
	if p.ContentTR.Projects == nil {
		p.ContentTR.Projects = []Project{}
	}
	if p.ContentTR.SkillGroups == nil {
		p.ContentTR.SkillGroups = []SkillGroup{}
	}
	if p.ContentEN.Experiences == nil {
		p.ContentEN.Experiences = []Experience{}
	}
	if p.ContentEN.Education == nil {
		p.ContentEN.Education = []Education{}
	}
	if p.ContentEN.Projects == nil {
		p.ContentEN.Projects = []Project{}
	}
	if p.ContentEN.SkillGroups == nil {
		p.ContentEN.SkillGroups = []SkillGroup{}
	}
	if p.ContentLanguages == nil {
		p.ContentLanguages = map[string]*LocalizedContent{}
	}
	for lang, content := range p.ContentLanguages {
		if content == nil {
			p.ContentLanguages[lang] = ptrLocalizedContent(EmptyLocalizedContent())
			continue
		}
		*p.ContentLanguages[lang] = ensureLocalizedContentSlices(*content)
	}
}

func ensureLocalizedContentSlices(content LocalizedContent) LocalizedContent {
	if content.Experiences == nil {
		content.Experiences = []Experience{}
	}
	if content.Education == nil {
		content.Education = []Education{}
	}
	if content.Projects == nil {
		content.Projects = []Project{}
	}
	if content.SkillGroups == nil {
		content.SkillGroups = []SkillGroup{}
	}
	return content
}

func (p *CVProfile) hasLegacyContent() bool {
	if strings.TrimSpace(p.Summary) != "" {
		return true
	}
	if len(p.Experiences) > 0 || len(p.Education) > 0 || len(p.Projects) > 0 || len(p.SkillGroups) > 0 {
		return true
	}
	if strings.TrimSpace(p.Personal.Title) != "" || strings.TrimSpace(p.Personal.Languages) != "" {
		return true
	}
	return false
}

func (p *CVProfile) MigrateFromLegacy() {
	if !p.hasLegacyContent() {
		return
	}
	if !p.ContentTR.IsEmpty() || !p.ContentEN.IsEmpty() {
		p.clearLegacyFields()
		return
	}

	p.ContentTR = LocalizedContent{
		Summary:           p.Summary,
		Experiences:       append([]Experience(nil), p.Experiences...),
		Education:         append([]Education(nil), p.Education...),
		Projects:          append([]Project(nil), p.Projects...),
		SkillGroups:       append([]SkillGroup(nil), p.SkillGroups...),
		PersonalTitle:     p.Personal.Title,
		PersonalLanguages: p.Personal.Languages,
	}
	p.clearLegacyFields()
}

func (p *CVProfile) clearLegacyFields() {
	p.Summary = ""
	p.Experiences = nil
	p.Education = nil
	p.Projects = nil
	p.SkillGroups = nil
	p.Personal.Title = ""
	p.Personal.Languages = ""
}

func (p *CVProfile) NormalizePhotoSettings() {
	legacy := p.PhotoSize <= 0 && p.PhotoBorderWidth == 0 && p.PhotoBorderColor == ""
	if p.PhotoSize <= 0 {
		p.PhotoSize = DefaultPhotoSize
	}
	if p.PhotoBorderColor == "" {
		p.PhotoBorderColor = DefaultPhotoBorderColor
	}
	if legacy || p.PhotoBorderWidth < 0 {
		p.PhotoBorderWidth = DefaultPhotoBorderWidth
	}
}

func (p *CVProfile) UnmarshalJSON(data []byte) error {
	type profileAlias CVProfile
	aux := (*profileAlias)(p)
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	p.MigrateFromLegacy()
	p.ensureContentSlices()
	return nil
}

type ProfileSummary struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type ProfileCopyOptions struct {
	CopyPhoto     bool     `json:"copyPhoto"`
	CopyPersonal  bool     `json:"copyPersonal"`
	CopyContentTR bool     `json:"copyContentTR"`
	CopyContentEN bool     `json:"copyContentEN"`
	Sections      []string `json:"sections"`
}

func (o ProfileCopyOptions) CopiesAllContent() bool {
	if len(o.Sections) == 0 {
		return o.CopyContentTR || o.CopyContentEN
	}
	for _, section := range o.Sections {
		if strings.EqualFold(strings.TrimSpace(section), "all_content") {
			return true
		}
	}
	return false
}

func (o ProfileCopyOptions) WantsSection(name string) bool {
	if o.CopiesAllContent() {
		return true
	}
	for _, section := range o.Sections {
		if strings.EqualFold(strings.TrimSpace(section), name) {
			return true
		}
	}
	return false
}
