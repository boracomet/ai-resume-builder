package handlers

import (
	"database/sql"
	"net/http"
	"os"
	"strings"

	"github.com/boracomet/ai-resume-builder/internal/ai"
	"github.com/boracomet/ai-resume-builder/internal/db"
	"github.com/boracomet/ai-resume-builder/internal/models"
	"github.com/gin-gonic/gin"
)

type AIHandler struct {
	repo *db.CVRepository
}

func NewAIHandler(repo *db.CVRepository) *AIHandler {
	return &AIHandler{repo: repo}
}

type aiTestRequest struct {
	APIKey string `json:"apiKey"`
	Model  string `json:"model"`
}

type aiChatRequest struct {
	APIKey            string       `json:"apiKey"`
	Model             string       `json:"model"`
	Messages          []ai.Message `json:"messages"`
	ProfileID         *int64       `json:"profileId"`
	Language          string       `json:"language"`
	TranslateProvider string       `json:"translateProvider"`
}

type aiApplyRequest struct {
	Action      string                   `json:"action"`
	ProfileID   *int64                   `json:"profileId"`
	ProfileName string                   `json:"profileName"`
	Language    string                   `json:"language"`
	Content     *models.LocalizedContent `json:"content"`
}

func (h *AIHandler) resolveAPIKey(c *gin.Context, fromBody string) string {
	return ai.ResolveAPIKey(fromBody, c.Query("apiKey"), c.GetHeader("X-OpenAI-API-Key"))
}

func (h *AIHandler) ListModels(c *gin.Context) {
	apiKey := h.resolveAPIKey(c, c.Query("apiKey"))
	if apiKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "OpenAI API anahtarı gerekli (OPENAI_API_KEY veya arayüzden girin)"})
		return
	}

	models, err := ai.ListModels(apiKey)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"models": models})
}

func (h *AIHandler) TestConnection(c *gin.Context) {
	var req aiTestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	apiKey := h.resolveAPIKey(c, req.APIKey)
	if apiKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "OpenAI API anahtarı gerekli (OPENAI_API_KEY veya arayüzden girin)"})
		return
	}

	model := strings.TrimSpace(req.Model)
	if model == "" {
		model = strings.TrimSpace(os.Getenv("OPENAI_DEFAULT_MODEL"))
	}
	if model == "" {
		model = "gpt-4o-mini"
	}

	result, err := ai.TestConnection(apiKey, model)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error(), "ok": false})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"ok":      true,
		"message": "Bağlantı başarılı",
		"model":   model,
		"usage":   result.Usage,
		"costUSD": result.CostUSD,
	})
}

func (h *AIHandler) buildToolContext(req aiChatRequest, apiKey string) ai.ToolContext {
	provider := strings.ToLower(strings.TrimSpace(req.TranslateProvider))
	if provider == "" {
		provider = "google"
	}
	language := models.NormalizeLangCode(req.Language)
	if language == "" {
		language = "tr"
	}
	return ai.ToolContext{
		ProfileID:             req.ProfileID,
		Language:              language,
		TranslateProvider:     provider,
		OpenAIAPIKey:          apiKey,
		OpenAIModel:           strings.TrimSpace(req.Model),
		GoogleTranslateAPIKey: strings.TrimSpace(os.Getenv("GOOGLE_TRANSLATE_API_KEY")),
	}
}

func (h *AIHandler) Chat(c *gin.Context) {
	var req aiChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	apiKey := h.resolveAPIKey(c, req.APIKey)
	if apiKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "OpenAI API anahtarı gerekli (OPENAI_API_KEY veya arayüzden girin)"})
		return
	}

	model := strings.TrimSpace(req.Model)
	if model == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "model seçimi gerekli"})
		return
	}
	if len(req.Messages) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "mesaj listesi boş olamaz"})
		return
	}

	var profile *models.CVProfile
	if req.ProfileID != nil && *req.ProfileID > 0 {
		loaded, err := h.repo.GetByID(*req.ProfileID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		profile = loaded
	}

	messages := []ai.Message{
		{Role: "system", Content: ai.SystemPrompt},
	}
	if profile != nil {
		messages = append(messages, ai.Message{
			Role:    "system",
			Content: ai.BuildProfileContext(profile),
		})
	}
	messages = append(messages, req.Messages...)

	toolCtx := h.buildToolContext(req, apiKey)
	toolState := &ai.ToolExecutionState{}
	executor := ai.NewToolExecutor(h.repo)

	chatResult, err := ai.RunToolChatLoop(apiKey, model, messages, executor, &toolCtx, toolState)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	response := gin.H{
		"message":        chatResult.Reply,
		"profileUpdated": toolState.ProfileUpdated,
		"actions":        toolState.Actions,
		"usage":          chatResult.Usage,
		"costUSD":        chatResult.CostUSD,
	}
	if toolState.ProfileID != nil {
		response["profileId"] = *toolState.ProfileID
	}
	if toolState.Language != "" {
		response["language"] = toolState.Language
	}
	if toolState.Profile != nil {
		response["profile"] = toolState.Profile
	}
	if toolState.ShowProfilePicker {
		response["showProfilePicker"] = true
		response["profilePickerMessage"] = toolState.ProfilePickerMessage
		if toolState.CopyOptions != nil {
			response["copyOptions"] = toolState.CopyOptions
		}
	}

	c.JSON(http.StatusOK, response)
}

func (h *AIHandler) Apply(c *gin.Context) {
	var req aiApplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	action := strings.TrimSpace(req.Action)
	if action != "create_profile" && action != "update_profile" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "geçersiz action: create_profile veya update_profile olmalı"})
		return
	}
	if req.Content == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "content gerekli"})
		return
	}

	language := models.NormalizeLangCode(req.Language)
	if language == "" {
		language = "tr"
	}

	if action == "create_profile" {
		profile := models.CVProfile{
			Name:     strings.TrimSpace(req.ProfileName),
			Language: language,
		}
		if profile.Name == "" {
			profile.Name = "AI CV"
		}
		if language == "en" {
			profile.ContentEN = *req.Content
			profile.ContentTR = models.EmptyLocalizedContent()
		} else {
			profile.ContentTR = *req.Content
			profile.ContentEN = models.EmptyLocalizedContent()
		}
		profile.Normalize()

		if err := h.repo.Create(&profile); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, profile)
		return
	}

	if req.ProfileID == nil || *req.ProfileID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "update_profile için profileId gerekli"})
		return
	}

	profile, err := h.repo.GetByID(*req.ProfileID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if profile == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "profil bulunamadı"})
		return
	}

	if name := strings.TrimSpace(req.ProfileName); name != "" {
		profile.Name = name
	}
	profile.Language = language
	if language == "en" {
		profile.ContentEN = mergeLocalizedContent(profile.ContentEN, *req.Content)
	} else {
		profile.ContentTR = mergeLocalizedContent(profile.ContentTR, *req.Content)
	}
	profile.Normalize()

	if err := h.repo.Update(profile); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "profil bulunamadı"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	updated, err := h.repo.GetByID(*req.ProfileID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updated)
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
