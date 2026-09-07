package ocr

import (
	"bytes"
	"fmt"
	"image/png"
	"strings"

	"github.com/gen2brain/go-fitz"
)

const minCharsPerPage = 100

// ExtractTextFromPDF reads a PDF and returns the best available plain text.
func ExtractTextFromPDF(pdfBytes []byte) (string, error) {
	text, _, _, err := OCRPDF(pdfBytes)
	return text, err
}

// OCRPDF extracts text from a PDF. It uses embedded text when sufficient,
// otherwise falls back to Tesseract OCR on rendered page images.
func OCRPDF(pdfBytes []byte) (text string, pages int, method string, err error) {
	doc, err := fitz.NewFromMemory(pdfBytes)
	if err != nil {
		return "", 0, "", fmt.Errorf("PDF açılamadı: %w", err)
	}
	defer doc.Close()

	pageCount := doc.NumPage()
	if pageCount == 0 {
		return "", 0, "", fmt.Errorf("PDF boş veya sayfa içermiyor")
	}

	embedded := extractEmbeddedText(doc, pageCount)
	avgChars := len(embedded) / pageCount
	if avgChars >= minCharsPerPage {
		return embedded, pageCount, "extract", nil
	}

	ocrText, ocrErr := ocrWithTesseract(doc, pageCount)
	if ocrErr != nil {
		if strings.TrimSpace(embedded) != "" {
			return embedded, pageCount, "extract", nil
		}
		return "", pageCount, "", fmt.Errorf("OCR başarısız: %w", ocrErr)
	}

	if strings.TrimSpace(ocrText) == "" {
		if strings.TrimSpace(embedded) != "" {
			return embedded, pageCount, "extract", nil
		}
		return "", pageCount, "", fmt.Errorf("PDF'den metin çıkarılamadı")
	}

	return ocrText, pageCount, "tesseract", nil
}

func extractEmbeddedText(doc *fitz.Document, pageCount int) string {
	var builder strings.Builder
	for i := 0; i < pageCount; i++ {
		pageText, err := doc.Text(i)
		if err != nil {
			continue
		}
		pageText = strings.TrimSpace(pageText)
		if pageText == "" {
			continue
		}
		if builder.Len() > 0 {
			builder.WriteString("\n\n")
		}
		builder.WriteString(pageText)
	}
	return strings.TrimSpace(builder.String())
}

func ocrWithTesseract(doc *fitz.Document, pageCount int) (string, error) {
	if !tesseractAvailable() {
		return "", fmt.Errorf(tesseractInstallHint)
	}

	var builder strings.Builder
	for i := 0; i < pageCount; i++ {
		img, err := doc.ImageDPI(i, 200)
		if err != nil {
			return "", fmt.Errorf("sayfa %d görüntüye dönüştürülemedi: %w", i+1, err)
		}

		var buf bytes.Buffer
		if err := png.Encode(&buf, img); err != nil {
			return "", fmt.Errorf("sayfa %d PNG'ye kodlanamadı: %w", i+1, err)
		}

		pageText, err := runTesseractCLI(buf.Bytes())
		if err != nil {
			return "", fmt.Errorf("sayfa %d OCR hatası: %w", i+1, err)
		}

		if pageText == "" {
			continue
		}
		if builder.Len() > 0 {
			builder.WriteString("\n\n")
		}
		builder.WriteString(pageText)
	}

	return strings.TrimSpace(builder.String()), nil
}
