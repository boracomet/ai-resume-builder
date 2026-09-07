package handlers

import (
	"io"
	"net/http"
	"strings"

	"github.com/boracomet/ai-resume-builder/internal/ocr"
	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
)

const maxOCRPDFSize = 10 << 20 // 10 MB

type OCRHandler struct{}

func NewOCRHandler() *OCRHandler {
	return &OCRHandler{}
}

func (h *OCRHandler) ProcessPDF(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		file, err = c.FormFile("pdf")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "PDF dosyası gerekli (alan adı: file veya pdf)"})
			return
		}
	}

	if file.Size > maxOCRPDFSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dosya boyutu 10 MB sınırını aşıyor"})
		return
	}

	if !strings.EqualFold(strings.TrimSpace(file.Header.Get("Content-Type")), "application/pdf") {
		name := strings.ToLower(file.Filename)
		if !strings.HasSuffix(name, ".pdf") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Yalnızca PDF dosyaları desteklenir"})
			return
		}
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer src.Close()

	data, err := io.ReadAll(io.LimitReader(src, maxOCRPDFSize+1))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(data) > maxOCRPDFSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dosya boyutu 10 MB sınırını aşıyor"})
		return
	}

	mime := mimetype.Detect(data)
	if mime != nil && mime.String() != "application/pdf" && !strings.HasSuffix(strings.ToLower(file.Filename), ".pdf") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz dosya türü; PDF bekleniyor"})
		return
	}

	text, pages, method, err := ocr.OCRPDF(data)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"text":   text,
		"pages":  pages,
		"chars":  len(text),
		"method": method,
		"file":   file.Filename,
	})
}
