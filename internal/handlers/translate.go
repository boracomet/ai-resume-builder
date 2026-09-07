package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/boracomet/ai-resume-builder/internal/db"
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
	SourceLang string `json:"sourceLang"`
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

	targetLang := translate.NormalizeLangCode(req.TargetLang)
	if !translate.ValidLangCode(targetLang) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "geçersiz hedef dil kodu (ISO 639-1)"})
		return
	}
	storageSupported := translate.StorageSupported(targetLang)

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

	sourceLang := translate.ResolveSourceLang(req.SourceLang, profile.Language, targetLang)
	if !translate.ValidLangCode(sourceLang) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "geçersiz kaynak dil kodu (ISO 639-1)"})
		return
	}
	if sourceLang == targetLang {
		c.JSON(http.StatusBadRequest, gin.H{"error": "kaynak ve hedef dil aynı olamaz"})
		return
	}

	var translatedContent interface{}
	if storageSupported {
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
	} else {
		content, err := translate.TranslateContent(profile, sourceLang, targetLang, func(texts []string, src, tgt string) ([]string, error) {
			switch provider {
			case "openai":
				translated, u, c, err := translate.TranslateTextsWithOpenAI(texts, src, tgt, apiKey, model)
				if err != nil {
					return nil, err
				}
				usage = u
				costUSD = c
				return translated, nil
			default:
				return translate.TranslateTexts(texts, src, tgt, apiKey)
			}
		})
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		translatedContent = content
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
	response["storageSupported"] = storageSupported
	response["targetLang"] = targetLang
	response["sourceLang"] = sourceLang
	if !storageSupported && translatedContent != nil {
		response["translatedContent"] = translatedContent
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
