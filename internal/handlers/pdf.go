package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/boracomet/ai-resume-builder/internal/models"
	"github.com/boracomet/ai-resume-builder/internal/pdf"
	"github.com/boracomet/ai-resume-builder/internal/render"
	"github.com/gin-gonic/gin"
)

type PDFHandler struct{}

func NewPDFHandler() *PDFHandler {
	return &PDFHandler{}
}

func (h *PDFHandler) Generate(c *gin.Context) {
	var profile models.CVProfile
	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
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

	pdfBytes, err := pdf.Generate(html)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("PDF oluşturulamadı: %v", err)})
		return
	}

	filename := "resume.pdf"
	if name := strings.TrimSpace(profile.Personal.Name); name != "" {
		slug := strings.ToLower(strings.ReplaceAll(name, " ", "-"))
		filename = slug + "-cv.pdf"
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}
