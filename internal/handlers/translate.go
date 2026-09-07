package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/boracomet/ai-resume-builder/internal/db"
	"github.com/boracomet/ai-resume-builder/internal/i18n"
	"github.com/boracomet/ai-resume-builder/internal/openaiusage"
	"github.com/boracomet/ai-resume-builder/internal/translate"
	"github.com/gin-gonic/gin"
)

type TranslateHandler struct {
	repo *db.CVRepository
}

func NewTranslateHandler(repo *db.CVRepository) *TranslateHandler {
	return &TranslateHandler{repo: repo}
}

type translateRequest struct {
	TargetLang string `json:"targetLang"`
	APIKey     string `json:"apiKey"`
	Provider   string `json:"provider"`
	Model      string `json:"model"`
}

func (h *TranslateHandler) TranslateProfile(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	var req translateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	targetLang := i18n.NormalizeLang(req.TargetLang)
	if targetLang != "en" && targetLang != "tr" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "desteklenen hedef diller: tr, en"})
		return
	}

	provider := strings.ToLower(strings.TrimSpace(req.Provider))
	if provider == "" {
		provider = "google"
	}

	apiKey := strings.TrimSpace(req.APIKey)
	model := strings.TrimSpace(req.Model)

	var translateErr error
	var usage openaiusage.Usage
	var costUSD float64
	switch provider {
	case "openai":
		if apiKey == "" {
			apiKey = strings.TrimSpace(os.Getenv("OPENAI_API_KEY"))
		}
		if apiKey == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "OpenAI API anahtarı gerekli (OPENAI_API_KEY veya arayüzden girin)"})
			return
		}
		if model == "" {
			model = strings.TrimSpace(os.Getenv("OPENAI_DEFAULT_MODEL"))
		}
		if model == "" {
			model = "gpt-4o-mini"
		}
	default:
		provider = "google"
		if apiKey == "" {
			apiKey = strings.TrimSpace(os.Getenv("GOOGLE_TRANSLATE_API_KEY"))
		}
		if apiKey == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Google Translate API anahtarı gerekli (GOOGLE_TRANSLATE_API_KEY veya arayüzden girin)"})
			return
		}
	}

	profile, err := h.repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if profile == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "profil bulunamadı"})
		return
	}

	sourceLang := "tr"
	if targetLang == "tr" {
		sourceLang = "en"
	}
	switch provider {
	case "openai":
		usage, costUSD, translateErr = translate.TranslateProfileWithOpenAI(profile, sourceLang, targetLang, apiKey, model)
	default:
		translateErr = translate.TranslateProfile(profile, sourceLang, targetLang, apiKey)
	}
	if translateErr != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": translateErr.Error()})
		return
	}

	if err := h.repo.Update(profile); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "profil bulunamadı"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	updated, err := h.repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := gin.H{}
	profileJSON, err := json.Marshal(updated)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := json.Unmarshal(profileJSON, &response); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	switch provider {
	case "openai":
		if usage.PromptTokens > 0 || usage.CompletionTokens > 0 {
			response["usage"] = usage
			response["costUSD"] = costUSD
		}
		response["provider"] = "openai"
	default:
		response["provider"] = "google"
		response["costUSD"] = 0
	}
	c.JSON(http.StatusOK, response)
}
