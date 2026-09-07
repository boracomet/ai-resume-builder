package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/boracomet/ai-resume-builder/internal/db"
	"github.com/boracomet/ai-resume-builder/internal/models"
	"github.com/gin-gonic/gin"
)

const (
	backupVersion       = "1"
	backupApp           = "ai-resume-builder"
	maxImportBodySize   = 5 * 1024 * 1024
	maxPhotoBase64Bytes = 2 * 1024 * 1024
)

type BackupHandler struct {
	repo *db.CVRepository
}

func NewBackupHandler(repo *db.CVRepository) *BackupHandler {
	return &BackupHandler{repo: repo}
}

type BackupPayload struct {
	Version    string              `json:"version"`
	App        string              `json:"app"`
	ExportedAt time.Time           `json:"exportedAt"`
	Profiles   []models.CVProfile  `json:"profiles"`
}

type ImportRequest struct {
	Mode     string             `json:"mode"`
	Profiles []models.CVProfile `json:"profiles"`
}

func (h *BackupHandler) ExportProfile(c *gin.Context) {
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

	payload := newBackupPayload([]models.CVProfile{*profile})
	filename := sanitizeFilename(profile.Name) + ".json"
	writeBackupJSON(c, payload, filename)
}

func (h *BackupHandler) ExportAll(c *gin.Context) {
	profiles, err := h.repo.ExportAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	payload := newBackupPayload(profiles)
	filename := fmt.Sprintf("ai-resume-builder-backup-%s.json", time.Now().UTC().Format("2006-01-02"))
	writeBackupJSON(c, payload, filename)
}

func (h *BackupHandler) Import(c *gin.Context) {
	body, err := readImportBody(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(body) > maxImportBodySize {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("dosya boyutu %d MB limitini aşıyor", maxImportBodySize/(1024*1024))})
		return
	}

	mode, profiles, err := parseImportPayload(body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(profiles) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "içe aktarılacak profil bulunamadı"})
		return
	}

	for i := range profiles {
		if err := validateProfileForImport(&profiles[i]); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("profil %d: %v", i+1, err)})
			return
		}
	}

	imported, err := h.repo.ImportProfiles(profiles, mode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  fmt.Sprintf("%d profil içe aktarıldı", len(imported)),
		"count":    len(imported),
		"profiles": imported,
	})
}

func newBackupPayload(profiles []models.CVProfile) BackupPayload {
	exportProfiles := make([]models.CVProfile, len(profiles))
	for i, p := range profiles {
		exportProfiles[i] = p
		exportProfiles[i].ID = 0
		exportProfiles[i].CreatedAt = time.Time{}
		exportProfiles[i].UpdatedAt = time.Time{}
	}

	return BackupPayload{
		Version:    backupVersion,
		App:        backupApp,
		ExportedAt: time.Now().UTC(),
		Profiles:   exportProfiles,
	}
}

func writeBackupJSON(c *gin.Context, payload BackupPayload, filename string) {
	data, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Data(http.StatusOK, "application/json; charset=utf-8", data)
}

func readImportBody(c *gin.Context) ([]byte, error) {
	file, err := c.FormFile("file")
	if err == nil {
		src, err := file.Open()
		if err != nil {
			return nil, err
		}
		defer src.Close()
		return io.ReadAll(io.LimitReader(src, maxImportBodySize+1))
	}

	return io.ReadAll(io.LimitReader(c.Request.Body, maxImportBodySize+1))
}

func parseImportPayload(body []byte) (string, []models.CVProfile, error) {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(body, &raw); err != nil {
		return "", nil, fmt.Errorf("geçersiz JSON: %w", err)
	}

	mode := "merge"
	if m, ok := raw["mode"]; ok {
		var modeStr string
		if err := json.Unmarshal(m, &modeStr); err == nil {
			mode = normalizeImportMode(modeStr)
		}
	}

	profilesRaw, ok := raw["profiles"]
	if !ok {
		return "", nil, fmt.Errorf("profiles alanı bulunamadı")
	}

	var profiles []models.CVProfile
	if err := json.Unmarshal(profilesRaw, &profiles); err != nil {
		return "", nil, fmt.Errorf("profiles parse hatası: %w", err)
	}
	if len(profiles) == 0 {
		return "", nil, fmt.Errorf("profiles alanı boş")
	}
	return mode, profiles, nil
}

func normalizeImportMode(mode string) string {
	mode = strings.ToLower(strings.TrimSpace(mode))
	if mode == "replace" {
		return "replace"
	}
	return "merge"
}

func validateProfileForImport(profile *models.CVProfile) error {
	profile.Normalize()
	if strings.TrimSpace(profile.Name) == "" {
		profile.Name = "Imported Profile"
	}
	if len(profile.PhotoBase64) > maxPhotoBase64Bytes {
		return fmt.Errorf("photoBase64 boyutu limiti aşıyor (max %d KB)", maxPhotoBase64Bytes/1024)
	}
	return nil
}

func sanitizeFilename(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return "profile"
	}
	replacer := strings.NewReplacer(
		"/", "-",
		"\\", "-",
		":", "-",
		"*", "-",
		"?", "-",
		"\"", "-",
		"<", "-",
		">", "-",
		"|", "-",
	)
	return replacer.Replace(name)
}
