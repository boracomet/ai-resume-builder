package handlers

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type SettingsHandler struct{}

func NewSettingsHandler() *SettingsHandler {
	return &SettingsHandler{}
}

type settingsStatusResponse struct {
	GoogleTranslateConfigured bool `json:"googleTranslateConfigured"`
	OpenAIConfigured          bool `json:"openaiConfigured"`
	GoogleGeminiConfigured    bool `json:"googleGeminiConfigured"`
	XiaomiConfigured          bool `json:"xiaomiConfigured"`
}

func envConfigured(key string) bool {
	return strings.TrimSpace(os.Getenv(key)) != ""
}

func (h *SettingsHandler) GetStatus(c *gin.Context) {
	c.JSON(http.StatusOK, settingsStatusResponse{
		GoogleTranslateConfigured: envConfigured("GOOGLE_TRANSLATE_API_KEY"),
		OpenAIConfigured:          envConfigured("OPENAI_API_KEY"),
		GoogleGeminiConfigured:    envConfigured("GOOGLE_GEMINI_API_KEY"),
		XiaomiConfigured:          envConfigured("XIAOMI_API_KEY"),
	})
}
