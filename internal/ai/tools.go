package ai

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/boracomet/ai-resume-builder/internal/db"
	"github.com/boracomet/ai-resume-builder/internal/models"
	"github.com/boracomet/ai-resume-builder/internal/openaiusage"
	"github.com/boracomet/ai-resume-builder/internal/translate"
)

const maxToolIterations = 8

type ToolContext struct {
	ProfileID             *int64
	Language              string
	TranslateProvider     string
	OpenAIAPIKey          string
	OpenAIModel           string
	GoogleTranslateAPIKey string
}

type ToolAction struct {
	Tool    string `json:"tool"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ToolExecutionState struct {
	ProfileUpdated       bool                       `json:"profileUpdated"`
	ProfileID            *int64                     `json:"profileId,omitempty"`
	Profile              *models.CVProfile          `json:"profile,omitempty"`
	Language             string                     `json:"language,omitempty"`
	Actions              []ToolAction               `json:"actions"`
	ShowProfilePicker    bool                       `json:"showProfilePicker"`
	ProfilePickerMessage string                     `json:"profilePickerMessage,omitempty"`
	CopyOptions          *models.ProfileCopyOptions `json:"copyOptions,omitempty"`
	LastUsage            openaiusage.Usage
	LastCostUSD          float64
}

type ToolExecutor struct {
	repo *db.CVRepository
}

func NewToolExecutor(repo *db.CVRepository) *ToolExecutor {
	return &ToolExecutor{repo: repo}
}

func ToolDefinitions() []ToolDefinition {
	return []ToolDefinition{
		{
			Type: "function",
			Function: ToolFunctionSchema{
				Name:        "create_profile",
				Description: "Yeni boş CV profili oluşturur. Kullanıcı 'yeni profil aç', 'yeni cv oluştur' dediğinde kullan.",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"name": map[string]interface{}{
							"type":        "string",
							"description": "Profil adı (opsiyonel)",
						},
						"language": map[string]interface{}{
							"type":        "string",
							"enum":        []string{"tr", "en"},
							"description": "Düzenleme dili (varsayılan: tr)",
						},
					},
				},
			},
		},
		{
			Type: "function",
			Function: ToolFunctionSchema{
				Name:        "duplicate_profile",
				Description: "Mevcut CV profilinin kopyasını oluşturur.",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"profileId": map[string]interface{}{
							"type":        "integer",
							"description": "Kopyalanacak profil ID (boşsa aktif profil)",
						},
					},
				},
			},
		},
		{
			Type: "function",
			Function: ToolFunctionSchema{
				Name:        "update_profile_content",
				Description: "Aktif dil slotundaki CV içeriğini günceller (özet, deneyimler, eğitim, beceriler vb.).",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"profileId": map[string]interface{}{
							"type":        "integer",
							"description": "Güncellenecek profil ID (boşsa aktif profil)",
						},
						"language": map[string]interface{}{
							"type":        "string",
							"enum":        []string{"tr", "en"},
							"description": "İçeriğin yazılacağı dil",
						},
						"profileName": map[string]interface{}{
							"type":        "string",
							"description": "Profil adı (opsiyonel)",
						},
						"content": localizedContentSchema(),
					},
					"required": []string{"content"},
				},
			},
		},
		{
			Type: "function",
			Function: ToolFunctionSchema{
				Name:        "translate_profile",
				Description: "Profili Türkçe↔İngilizce çevirir ve kaydeder. 'ingilizce çevir', 'türkçeye çevir' isteklerinde kullan.",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"profileId": map[string]interface{}{
							"type":        "integer",
							"description": "Çevrilecek profil ID (boşsa aktif profil)",
						},
						"targetLang": map[string]interface{}{
							"type":        "string",
							"enum":        []string{"tr", "en"},
							"description": "Hedef dil",
						},
					},
					"required": []string{"targetLang"},
				},
			},
		},
		{
			Type: "function",
			Function: ToolFunctionSchema{
				Name:        "switch_language",
				Description: "Düzenleme dilini tr veya en olarak değiştirir.",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"profileId": map[string]interface{}{
							"type":        "integer",
							"description": "Profil ID (boşsa aktif profil)",
						},
						"language": map[string]interface{}{
							"type":        "string",
							"enum":        []string{"tr", "en"},
							"description": "Yeni düzenleme dili",
						},
					},
					"required": []string{"language"},
				},
			},
		},
		{
			Type: "function",
			Function: ToolFunctionSchema{
				Name:        "get_current_profile",
				Description: "Mevcut profil verisini okur; bağlam için kullan.",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"profileId": map[string]interface{}{
							"type":        "integer",
							"description": "Profil ID (boşsa aktif profil)",
						},
					},
				},
			},
		},
		{
			Type: "function",
			Function: ToolFunctionSchema{
				Name:        "apply_to_profile",
				Description: "Üretilen CV içeriğini profile uygular; yeni profil oluşturma veya mevcut profili güncelleme.",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"action": map[string]interface{}{
							"type":        "string",
							"enum":        []string{"create_profile", "update_profile"},
							"description": "create_profile veya update_profile",
						},
						"profileId": map[string]interface{}{
							"type":        "integer",
							"description": "update_profile için profil ID",
						},
						"profileName": map[string]interface{}{
							"type":        "string",
							"description": "Profil adı",
						},
						"language": map[string]interface{}{
							"type":        "string",
							"enum":        []string{"tr", "en"},
							"description": "İçerik dili",
						},
						"content": localizedContentSchema(),
					},
					"required": []string{"action", "content"},
				},
			},
		},
		{
			Type: "function",
			Function: ToolFunctionSchema{
				Name:        "request_profile_selection",
				Description: "Kullanıcı başka bir profilden veri almak istediğinde ancak hangi profil olduğunu belirtmediğinde veya belirsiz olduğunda çağır. Sohbette profil seçim butonları gösterilir.",
				Parameters:  profileCopyOptionsSchema(true),
			},
		},
		{
			Type: "function",
			Function: ToolFunctionSchema{
				Name:        "copy_from_profile",
				Description: "Belirli bir kaynak profilden mevcut profile fotoğraf, kişisel bilgiler ve/veya içerik kopyalar. Kaynak profil adı veya ID biliniyorsa kullan.",
				Parameters:  profileCopyOptionsSchema(false),
			},
		},
	}
}

func profileCopyOptionsSchema(includeMessage bool) map[string]interface{} {
	properties := map[string]interface{}{
		"profileId": map[string]interface{}{
			"type":        "integer",
			"description": "Hedef profil ID (boşsa aktif profil)",
		},
		"sourceProfileId": map[string]interface{}{
			"type":        "integer",
			"description": "Kaynak profil ID",
		},
		"sourceProfileName": map[string]interface{}{
			"type":        "string",
			"description": "Kaynak profil adı (ID bilinmiyorsa)",
		},
		"copyPhoto": map[string]interface{}{
			"type":        "boolean",
			"description": "Profil fotoğrafını kopyala",
		},
		"copyPersonal": map[string]interface{}{
			"type":        "boolean",
			"description": "Kişisel bilgileri (ad, e-posta, telefon vb.) kopyala",
		},
		"copyContentTR": map[string]interface{}{
			"type":        "boolean",
			"description": "Türkçe içeriği kopyala",
		},
		"copyContentEN": map[string]interface{}{
			"type":        "boolean",
			"description": "İngilizce içeriği kopyala",
		},
		"sections": map[string]interface{}{
			"type": "array",
			"items": map[string]interface{}{
				"type": "string",
				"enum": []string{"all_content", "summary", "experience", "education", "projects", "skills"},
			},
			"description": "Kopyalanacak bölümler; all_content tüm içeriği kopyalar",
		},
	}
	if includeMessage {
		properties["message"] = map[string]interface{}{
			"type":        "string",
			"description": "Profil seçim ekranında gösterilecek mesaj",
		}
	}
	return map[string]interface{}{
		"type":       "object",
		"properties": properties,
	}
}

func localizedContentSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"summary":           map[string]interface{}{"type": "string"},
			"personalTitle":     map[string]interface{}{"type": "string"},
			"personalLanguages": map[string]interface{}{"type": "string"},
			"experiences": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"title":       map[string]interface{}{"type": "string"},
						"company":     map[string]interface{}{"type": "string"},
						"startDate":   map[string]interface{}{"type": "string"},
						"endDate":     map[string]interface{}{"type": "string"},
						"duration":    map[string]interface{}{"type": "string"},
						"location":    map[string]interface{}{"type": "string"},
						"description": map[string]interface{}{"type": "string"},
						"highlights":  map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "string"}},
					},
				},
			},
			"education": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"degree": map[string]interface{}{
							"type":        "string",
							"description": "Program veya derece adı (ör. Grafik Tasarım, Web Tasarımı ve Kodlama)",
						},
						"institution": map[string]interface{}{
							"type":        "string",
							"description": "Okul veya üniversite adı (ör. Anadolu Üniversitesi). degree alanına yazma.",
						},
						"startDate": map[string]interface{}{
							"type":        "string",
							"description": "Başlangıç yılı veya tarih (ör. 2022)",
						},
						"endDate": map[string]interface{}{
							"type":        "string",
							"description": "Bitiş yılı veya 'Devam Ediyor' / 'Ongoing'",
						},
						"description": map[string]interface{}{
							"type":        "string",
							"description": "Ek açıklama metni; program adını buraya yazma",
						},
					},
				},
			},
			"projects": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"name":        map[string]interface{}{"type": "string"},
						"url":         map[string]interface{}{"type": "string"},
						"description": map[string]interface{}{"type": "string"},
					},
				},
			},
			"skillGroups": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"category": map[string]interface{}{"type": "string"},
						"skills":   map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "string"}},
					},
				},
			},
		},
	}
}

func (e *ToolExecutor) Execute(toolName string, args map[string]interface{}, ctx *ToolContext, state *ToolExecutionState) (string, error) {
	toolName = strings.TrimSpace(toolName)
	var result interface{}
	var err error

	switch toolName {
	case "create_profile":
		result, err = e.createProfile(args, ctx, state)
	case "duplicate_profile":
		result, err = e.duplicateProfile(args, ctx, state)
	case "update_profile_content":
		result, err = e.updateProfileContent(args, ctx, state)
	case "translate_profile":
		result, err = e.translateProfile(args, ctx, state)
	case "switch_language":
		result, err = e.switchLanguage(args, ctx, state)
	case "get_current_profile":
		result, err = e.getCurrentProfile(args, ctx)
	case "apply_to_profile":
		result, err = e.applyToProfile(args, ctx, state)
	case "request_profile_selection":
		result, err = e.requestProfileSelection(args, ctx, state)
	case "copy_from_profile":
		result, err = e.copyFromProfile(args, ctx, state)
	default:
		err = fmt.Errorf("bilinmeyen araç: %s", toolName)
	}

	status := "success"
	message := toolSuccessMessage(toolName, result)
	if err != nil {
		status = "error"
		message = err.Error()
		result = map[string]interface{}{"error": err.Error()}
	}

	state.Actions = append(state.Actions, ToolAction{
		Tool:    toolName,
		Status:  status,
		Message: message,
	})

	payload, marshalErr := json.Marshal(result)
	if marshalErr != nil {
		return "", marshalErr
	}
	if err != nil {
		return string(payload), err
	}
	return string(payload), nil
}

func toolSuccessMessage(toolName string, result interface{}) string {
	switch toolName {
	case "create_profile", "duplicate_profile", "apply_to_profile":
		if m, ok := result.(map[string]interface{}); ok {
			if name, ok := m["name"].(string); ok && name != "" {
				return fmt.Sprintf("Profil: %s", name)
			}
		}
	case "translate_profile":
		if m, ok := result.(map[string]interface{}); ok {
			if lang, ok := m["language"].(string); ok {
				return fmt.Sprintf("Dil: %s", lang)
			}
		}
	case "switch_language":
		if m, ok := result.(map[string]interface{}); ok {
			if lang, ok := m["language"].(string); ok {
				return fmt.Sprintf("Düzenleme dili: %s", lang)
			}
		}
	case "copy_from_profile":
		if m, ok := result.(map[string]interface{}); ok {
			if name, ok := m["sourceName"].(string); ok && name != "" {
				return fmt.Sprintf("Kaynak profil: %s", name)
			}
		}
	}
	return "Tamamlandı"
}

func (e *ToolExecutor) resolveProfileID(args map[string]interface{}, ctx *ToolContext) (int64, error) {
	if id, ok := argInt64(args, "profileId"); ok && id > 0 {
		return id, nil
	}
	if ctx.ProfileID != nil && *ctx.ProfileID > 0 {
		return *ctx.ProfileID, nil
	}
	return 0, fmt.Errorf("profil seçili değil; profileId belirtin veya önce profil oluşturun")
}

func (e *ToolExecutor) loadProfile(id int64) (*models.CVProfile, error) {
	profile, err := e.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if profile == nil {
		return nil, fmt.Errorf("profil bulunamadı: %d", id)
	}
	return profile, nil
}

func (e *ToolExecutor) markProfileUpdated(state *ToolExecutionState, profile *models.CVProfile) {
	state.ProfileUpdated = true
	id := profile.ID
	state.ProfileID = &id
	state.Profile = profile
	state.Language = profile.Language
	ctxLang := profile.Language
	_ = ctxLang
}

func (e *ToolExecutor) createProfile(args map[string]interface{}, ctx *ToolContext, state *ToolExecutionState) (interface{}, error) {
	name := strings.TrimSpace(argString(args, "name"))
	if name == "" {
		name = "Yeni CV"
	}
	language := models.NormalizeLangCode(argString(args, "language"))
	if language == "" {
		language = models.NormalizeLangCode(ctx.Language)
	}
	if language == "" {
		language = "tr"
	}

	profile := models.CVProfile{
		Name:      name,
		Language:  language,
		ContentTR: models.EmptyLocalizedContent(),
		ContentEN: models.EmptyLocalizedContent(),
	}
	profile.Normalize()

	if err := e.repo.Create(&profile); err != nil {
		return nil, err
	}

	e.markProfileUpdated(state, &profile)
	ctx.ProfileID = &profile.ID
	ctx.Language = profile.Language

	return map[string]interface{}{
		"profileId": profile.ID,
		"name":      profile.Name,
		"language":  profile.Language,
	}, nil
}

func (e *ToolExecutor) duplicateProfile(args map[string]interface{}, ctx *ToolContext, state *ToolExecutionState) (interface{}, error) {
	id, err := e.resolveProfileID(args, ctx)
	if err != nil {
		return nil, err
	}

	profile, err := e.repo.Duplicate(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("profil bulunamadı")
		}
		return nil, err
	}

	e.markProfileUpdated(state, profile)
	ctx.ProfileID = &profile.ID
	ctx.Language = profile.Language

	return map[string]interface{}{
		"profileId": profile.ID,
		"name":      profile.Name,
		"language":  profile.Language,
	}, nil
}

func (e *ToolExecutor) updateProfileContent(args map[string]interface{}, ctx *ToolContext, state *ToolExecutionState) (interface{}, error) {
	id, err := e.resolveProfileID(args, ctx)
	if err != nil {
		return nil, err
	}

	profile, err := e.loadProfile(id)
	if err != nil {
		return nil, err
	}

	content, err := parseLocalizedContent(args["content"])
	if err != nil {
		return nil, err
	}

	language := models.NormalizeLangCode(argString(args, "language"))
	if language == "" {
		language = models.NormalizeLangCode(ctx.Language)
	}
	if language == "" {
		language = profile.Language
	}

	if name := strings.TrimSpace(argString(args, "profileName")); name != "" {
		profile.Name = name
	}
	profile.Language = language

	existing := profile.ContentForLang(language)
	profile.SetContentForLang(language, mergeLocalizedContent(*existing, content))
	profile.Normalize()

	if err := e.repo.Update(profile); err != nil {
		return nil, err
	}

	updated, err := e.loadProfile(id)
	if err != nil {
		return nil, err
	}
	e.markProfileUpdated(state, updated)
	ctx.ProfileID = &updated.ID
	ctx.Language = updated.Language

	return map[string]interface{}{
		"profileId": updated.ID,
		"name":      updated.Name,
		"language":  updated.Language,
	}, nil
}

func (e *ToolExecutor) translateProfile(args map[string]interface{}, ctx *ToolContext, state *ToolExecutionState) (interface{}, error) {
	id, err := e.resolveProfileID(args, ctx)
	if err != nil {
		return nil, err
	}

	targetLang := translate.NormalizeLangCode(argString(args, "targetLang"))
	if !translate.ValidLangCode(targetLang) {
		return nil, fmt.Errorf("geçersiz hedef dil kodu (ISO 639-1)")
	}

	profile, err := e.loadProfile(id)
	if err != nil {
		return nil, err
	}

	sourceLang := translate.ResolveSourceLang(argString(args, "sourceLang"), profile.Language, targetLang)

	provider := strings.ToLower(strings.TrimSpace(ctx.TranslateProvider))
	if provider == "" {
		provider = "google"
	}

	switch provider {
	case "openai":
		apiKey := strings.TrimSpace(ctx.OpenAIAPIKey)
		if apiKey == "" {
			return nil, fmt.Errorf("OpenAI API anahtarı gerekli")
		}
		model := strings.TrimSpace(ctx.OpenAIModel)
		if model == "" {
			model = strings.TrimSpace(os.Getenv("OPENAI_DEFAULT_MODEL"))
		}
		if model == "" {
			model = "gpt-4o-mini"
		}
		usage, cost, err := translate.TranslateProfileWithOpenAI(profile, sourceLang, targetLang, apiKey, model)
		if err != nil {
			return nil, err
		}
		state.LastUsage = usage
		state.LastCostUSD = cost
	default:
		apiKey := strings.TrimSpace(ctx.GoogleTranslateAPIKey)
		if apiKey == "" {
			return nil, fmt.Errorf("Google Translate API anahtarı gerekli")
		}
		if err := translate.TranslateProfile(profile, sourceLang, targetLang, apiKey); err != nil {
			return nil, err
		}
	}

	if err := e.repo.Update(profile); err != nil {
		return nil, err
	}

	updated, err := e.loadProfile(id)
	if err != nil {
		return nil, err
	}
	e.markProfileUpdated(state, updated)
	ctx.ProfileID = &updated.ID
	ctx.Language = updated.Language
	state.Language = updated.Language

	return map[string]interface{}{
		"profileId": updated.ID,
		"name":      updated.Name,
		"language":  updated.Language,
		"targetLang": targetLang,
	}, nil
}

func (e *ToolExecutor) switchLanguage(args map[string]interface{}, ctx *ToolContext, state *ToolExecutionState) (interface{}, error) {
	language := models.NormalizeLangCode(argString(args, "language"))
	if !models.ValidLangCode(language) {
		return nil, fmt.Errorf("geçersiz dil kodu (ISO 639-1)")
	}

	id, err := e.resolveProfileID(args, ctx)
	if err != nil {
		return nil, err
	}

	profile, err := e.loadProfile(id)
	if err != nil {
		return nil, err
	}

	profile.Language = language
	profile.Normalize()

	if err := e.repo.Update(profile); err != nil {
		return nil, err
	}

	updated, err := e.loadProfile(id)
	if err != nil {
		return nil, err
	}
	e.markProfileUpdated(state, updated)
	ctx.Language = language
	state.Language = language

	return map[string]interface{}{
		"profileId": updated.ID,
		"language":  updated.Language,
	}, nil
}

func (e *ToolExecutor) getCurrentProfile(args map[string]interface{}, ctx *ToolContext) (interface{}, error) {
	id, err := e.resolveProfileID(args, ctx)
	if err != nil {
		return nil, err
	}

	profile, err := e.loadProfile(id)
	if err != nil {
		return nil, err
	}

	profile.Normalize()
	content := profile.ActiveContent()
	return map[string]interface{}{
		"id":       profile.ID,
		"name":     profile.Name,
		"language": profile.Language,
		"personal": profile.Personal,
		"content":  content,
	}, nil
}

func (e *ToolExecutor) applyToProfile(args map[string]interface{}, ctx *ToolContext, state *ToolExecutionState) (interface{}, error) {
	action := strings.TrimSpace(argString(args, "action"))
	content, err := parseLocalizedContent(args["content"])
	if err != nil {
		return nil, err
	}

	language := models.NormalizeLangCode(argString(args, "language"))
	if language == "" {
		language = models.NormalizeLangCode(ctx.Language)
	}
	if language == "" {
		language = "tr"
	}

	switch action {
	case "create_profile":
		name := strings.TrimSpace(argString(args, "profileName"))
		if name == "" {
			name = "AI CV"
		}
		profile := models.CVProfile{
			Name:     name,
			Language: language,
		}
		if language == "en" {
			profile.ContentEN = content
			profile.ContentTR = models.EmptyLocalizedContent()
		} else {
			profile.ContentTR = content
			profile.ContentEN = models.EmptyLocalizedContent()
		}
		profile.Normalize()

		if err := e.repo.Create(&profile); err != nil {
			return nil, err
		}
		e.markProfileUpdated(state, &profile)
		ctx.ProfileID = &profile.ID
		ctx.Language = profile.Language

		return map[string]interface{}{
			"profileId": profile.ID,
			"name":      profile.Name,
			"language":  profile.Language,
			"action":    action,
		}, nil

	case "update_profile":
		id, err := e.resolveProfileID(args, ctx)
		if err != nil {
			return nil, err
		}

		profile, err := e.loadProfile(id)
		if err != nil {
			return nil, err
		}

		if name := strings.TrimSpace(argString(args, "profileName")); name != "" {
			profile.Name = name
		}
		profile.Language = language
		if language == "en" {
			profile.ContentEN = mergeLocalizedContent(profile.ContentEN, content)
		} else {
			profile.ContentTR = mergeLocalizedContent(profile.ContentTR, content)
		}
		profile.Normalize()

		if err := e.repo.Update(profile); err != nil {
			return nil, err
		}

		updated, err := e.loadProfile(id)
		if err != nil {
			return nil, err
		}
		e.markProfileUpdated(state, updated)
		ctx.ProfileID = &updated.ID
		ctx.Language = updated.Language

		return map[string]interface{}{
			"profileId": updated.ID,
			"name":      updated.Name,
			"language":  updated.Language,
			"action":    action,
		}, nil

	default:
		return nil, fmt.Errorf("geçersiz action: create_profile veya update_profile olmalı")
	}
}

func (e *ToolExecutor) requestProfileSelection(args map[string]interface{}, ctx *ToolContext, state *ToolExecutionState) (interface{}, error) {
	opts := parseCopyOptions(args)
	message := strings.TrimSpace(argString(args, "message"))
	if message == "" {
		message = "Hangi profilden almak istiyorsunuz?"
	}

	state.ShowProfilePicker = true
	state.ProfilePickerMessage = message
	state.CopyOptions = &opts

	return map[string]interface{}{
		"showProfilePicker": true,
		"message":           message,
		"copyOptions":       opts,
	}, nil
}

func (e *ToolExecutor) copyFromProfile(args map[string]interface{}, ctx *ToolContext, state *ToolExecutionState) (interface{}, error) {
	targetID, err := e.resolveProfileID(args, ctx)
	if err != nil {
		return nil, err
	}

	opts := parseCopyOptions(args)
	sourceID, hasSourceID := argInt64(args, "sourceProfileId")
	if !hasSourceID || sourceID <= 0 {
		sourceName := strings.TrimSpace(argString(args, "sourceProfileName"))
		if sourceName == "" {
			return e.requestProfileSelection(args, ctx, state)
		}

		matches, err := e.findProfilesByName(sourceName)
		if err != nil {
			return nil, err
		}
		if len(matches) == 0 {
			return nil, fmt.Errorf("'%s' adına uygun profil bulunamadı", sourceName)
		}
		if len(matches) > 1 {
			message := fmt.Sprintf("'%s' için birden fazla profil var. Hangisini seçmek istiyorsunuz?", sourceName)
			args["message"] = message
			return e.requestProfileSelection(args, ctx, state)
		}
		sourceID = matches[0].ID
	}

	if sourceID == targetID {
		return nil, fmt.Errorf("kaynak ve hedef profil aynı olamaz")
	}

	updated, err := e.repo.CopyFrom(targetID, sourceID, opts)
	if err != nil {
		return nil, err
	}

	source, err := e.loadProfile(sourceID)
	if err != nil {
		return nil, err
	}

	e.markProfileUpdated(state, updated)
	ctx.ProfileID = &updated.ID
	ctx.Language = updated.Language

	return map[string]interface{}{
		"profileId":  updated.ID,
		"sourceId":   sourceID,
		"sourceName": source.Name,
		"copied":     opts,
	}, nil
}

func (e *ToolExecutor) findProfilesByName(name string) ([]models.ProfileSummary, error) {
	name = strings.TrimSpace(strings.ToLower(name))
	if name == "" {
		return nil, nil
	}

	profiles, err := e.repo.List()
	if err != nil {
		return nil, err
	}

	var matches []models.ProfileSummary
	for _, profile := range profiles {
		profileName := strings.TrimSpace(strings.ToLower(profile.Name))
		if profileName == name || strings.Contains(profileName, name) {
			matches = append(matches, profile)
		}
	}
	return matches, nil
}

func parseCopyOptions(args map[string]interface{}) models.ProfileCopyOptions {
	opts := models.ProfileCopyOptions{
		CopyPhoto:     argBool(args, "copyPhoto"),
		CopyPersonal:  argBool(args, "copyPersonal"),
		CopyContentTR: argBool(args, "copyContentTR"),
		CopyContentEN: argBool(args, "copyContentEN"),
	}
	if raw, ok := args["sections"]; ok && raw != nil {
		switch values := raw.(type) {
		case []interface{}:
			for _, value := range values {
				section := strings.TrimSpace(fmt.Sprint(value))
				if section != "" {
					opts.Sections = append(opts.Sections, section)
				}
			}
		case []string:
			opts.Sections = append(opts.Sections, values...)
		}
	}
	return opts
}

func argBool(args map[string]interface{}, key string) bool {
	value, ok := args[key]
	if !ok || value == nil {
		return false
	}
	switch v := value.(type) {
	case bool:
		return v
	case string:
		return strings.EqualFold(strings.TrimSpace(v), "true")
	default:
		return false
	}
}

func mergeLocalizedContent(existing models.LocalizedContent, patch models.LocalizedContent) models.LocalizedContent {
	if strings.TrimSpace(patch.Summary) != "" {
		existing.Summary = patch.Summary
	}
	if strings.TrimSpace(patch.PersonalTitle) != "" {
		existing.PersonalTitle = patch.PersonalTitle
	}
	if strings.TrimSpace(patch.PersonalLanguages) != "" {
		existing.PersonalLanguages = patch.PersonalLanguages
	}
	if len(patch.Experiences) > 0 {
		existing.Experiences = patch.Experiences
	}
	if len(patch.Education) > 0 {
		existing.Education = patch.Education
	}
	if len(patch.Projects) > 0 {
		existing.Projects = patch.Projects
	}
	if len(patch.SkillGroups) > 0 {
		existing.SkillGroups = patch.SkillGroups
	}
	return existing
}

func parseLocalizedContent(raw interface{}) (models.LocalizedContent, error) {
	if raw == nil {
		return models.LocalizedContent{}, fmt.Errorf("content gerekli")
	}
	data, err := json.Marshal(raw)
	if err != nil {
		return models.LocalizedContent{}, err
	}
	var content models.LocalizedContent
	if err := json.Unmarshal(data, &content); err != nil {
		return models.LocalizedContent{}, err
	}
	return content, nil
}

func argString(args map[string]interface{}, key string) string {
	value, ok := args[key]
	if !ok || value == nil {
		return ""
	}
	switch v := value.(type) {
	case string:
		return v
	case float64:
		return fmt.Sprintf("%.0f", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func argInt64(args map[string]interface{}, key string) (int64, bool) {
	value, ok := args[key]
	if !ok || value == nil {
		return 0, false
	}
	switch v := value.(type) {
	case float64:
		return int64(v), true
	case int:
		return int64(v), true
	case int64:
		return v, true
	case json.Number:
		parsed, err := v.Int64()
		return parsed, err == nil
	default:
		return 0, false
	}
}
