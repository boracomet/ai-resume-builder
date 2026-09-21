package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/boracomet/ai-resume-builder/internal/db"
	"github.com/boracomet/ai-resume-builder/internal/models"
	"github.com/boracomet/ai-resume-builder/internal/render"
	"github.com/gin-gonic/gin"
)

type PreviewHandler struct {
	repo *db.CVRepository
}

func NewPreviewHandler(repo *db.CVRepository) *PreviewHandler {
	return &PreviewHandler{repo: repo}
}

// previewRequest embeds CVProfile. Because CVProfile has a custom UnmarshalJSON,
// encoding/json would otherwise call that method with the entire payload and leave
// sibling fields (UseStoredPhoto) unbound. Unmarshal both parts explicitly.
type previewRequest struct {
	models.CVProfile
	UseStoredPhoto bool `json:"useStoredPhoto"`
}

func (r *previewRequest) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &r.CVProfile); err != nil {
		return err
	}
	var extra struct {
		UseStoredPhoto bool `json:"useStoredPhoto"`
	}
	if err := json.Unmarshal(data, &extra); err != nil {
		return err
	}
	r.UseStoredPhoto = extra.UseStoredPhoto
	return nil
}

func (h *PreviewHandler) Preview(c *gin.Context) {
	var req previewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profile := req.CVProfile
	if req.UseStoredPhoto && profile.ID > 0 && h.repo != nil {
		stored, err := h.repo.GetByID(profile.ID)
		if err == nil && stored != nil {
			profile.PhotoBase64 = stored.PhotoBase64
		}
	}

	renderer, err := render.Get()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	html, err := renderer.Render(&profile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Cache-Control", "no-store")
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}
