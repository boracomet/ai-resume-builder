package handlers

import (
	"database/sql"
	"encoding/base64"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/boradev/bora-cv/internal/db"
	"github.com/boradev/bora-cv/internal/models"
	"github.com/gin-gonic/gin"
)

type CVHandler struct {
	repo *db.CVRepository
}

func NewCVHandler(repo *db.CVRepository) *CVHandler {
	return &CVHandler{repo: repo}
}

func (h *CVHandler) ListProfiles(c *gin.Context) {
	profiles, err := h.repo.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, profiles)
}

func (h *CVHandler) GetProfile(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
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
	c.JSON(http.StatusOK, profile)
}

func (h *CVHandler) CreateProfile(c *gin.Context) {
	var profile models.CVProfile
	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if strings.TrimSpace(profile.Name) == "" {
		profile.Name = "Yeni CV"
	}
	profile.Normalize()

	if err := h.repo.Create(&profile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, profile)
}

func (h *CVHandler) UpdateProfile(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	var profile models.CVProfile
	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	profile.ID = id
	profile.Normalize()

	if err := h.repo.Update(&profile); err != nil {
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

func (h *CVHandler) DeleteProfile(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	if err := h.repo.Delete(id); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "profil bulunamadı"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "profil silindi"})
}

func (h *CVHandler) UploadPhoto(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	file, err := c.FormFile("photo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "fotoğraf dosyası gerekli"})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer src.Close()

	data, err := io.ReadAll(src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "image/jpeg"
	}

	encoded := base64.StdEncoding.EncodeToString(data)
	photoURI := "data:" + contentType + ";base64," + encoded

	profile, err := h.repo.UpdatePhoto(id, photoURI)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "profil bulunamadı"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, profile)
}

func parseID(c *gin.Context) (int64, bool) {
	id, err := parseInt64(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "geçersiz profil id"})
		return 0, false
	}
	return id, true
}

func parseInt64(value string) (int64, error) {
	return strconv.ParseInt(value, 10, 64)
}
