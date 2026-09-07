package ocr

import (
	"os"
	"strings"
	"testing"

	pdfgen "github.com/boracomet/ai-resume-builder/internal/pdf"
)

func TestOCRPDF_InvalidInput(t *testing.T) {
	_, _, _, err := OCRPDF([]byte("not a pdf"))
	if err == nil {
		t.Fatal("expected error for invalid PDF")
	}
}

func TestOCRPDF_EmptyPDF(t *testing.T) {
	_, _, _, err := OCRPDF([]byte{0x25, 0x50, 0x44, 0x46})
	if err == nil {
		t.Fatal("expected error for corrupt PDF")
	}
}

func TestOCRPDF_SampleEmbeddedText(t *testing.T) {
	html := `<!DOCTYPE html><html><body><h1>Bora Ata Turkoglu</h1><p>Senior Software Engineer with extensive experience in Go, distributed systems, and cloud infrastructure. Led multiple teams and delivered production systems at scale.</p></body></html>`
	data, err := pdfgen.Generate(html)
	if err != nil {
		t.Skipf("pdf generation unavailable: %v", err)
	}

	text, pages, method, err := OCRPDF(data)
	if err != nil {
		t.Fatalf("OCRPDF: %v", err)
	}
	if pages < 1 {
		t.Fatalf("expected at least 1 page, got %d", pages)
	}
	if method != "extract" {
		t.Fatalf("expected extract method, got %s", method)
	}
	if !strings.Contains(strings.ToLower(text), "software") {
		t.Fatalf("expected embedded text, got %q", text)
	}
}

func TestOCRPDF_FromEnv(t *testing.T) {
	path := os.Getenv("OCR_TEST_PDF")
	if path == "" {
		t.Skip("OCR_TEST_PDF not set")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read test pdf: %v", err)
	}

	text, pages, method, err := OCRPDF(data)
	if err != nil {
		t.Fatalf("OCRPDF: %v", err)
	}
	if pages <= 0 {
		t.Fatalf("expected pages > 0, got %d", pages)
	}
	if strings.TrimSpace(text) == "" {
		t.Fatal("expected non-empty text")
	}
	t.Logf("method=%s pages=%d chars=%d preview=%q", method, pages, len(text), truncate(text, 120))
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
