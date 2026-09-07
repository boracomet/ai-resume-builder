package handlers

import (
	"database/sql"
	"net/http"
	"os"
	"strings"

	"github.com/boradev/bora-cv/internal/db"
	"github.com/boradev/bora-cv/internal/i18n"
	"github.com/boradev/bora-cv/internal/translate"
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

	apiKey := strings.TrimSpace(req.APIKey)
	if apiKey == "" {
		apiKey = strings.TrimSpace(os.Getenv("GOOGLE_TRANSLATE_API_KEY"))
	}
	if apiKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Google Translate API anahtarı gerekli (GOOGLE_TRANSLATE_API_KEY veya arayüzden girin)"})
		return
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
	if err := translate.TranslateProfile(profile, sourceLang, targetLang, apiKey); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
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
	c.JSON(http.StatusOK, updated)
}
