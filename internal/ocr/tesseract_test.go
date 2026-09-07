package ocr

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"os"
	"os/exec"
	"testing"
)

func TestTesseractAvailable(t *testing.T) {
	available := tesseractAvailable()
	_, err := exec.LookPath("tesseract")
	if available != (err == nil) {
		t.Fatalf("tesseractAvailable=%v but LookPath err=%v", available, err)
	}
}

func TestRunTesseractCLI(t *testing.T) {
	if !tesseractAvailable() {
		t.Skip("tesseract not installed")
	}

	img := image.NewRGBA(image.Rect(0, 0, 200, 50))
	for x := 0; x < 200; x++ {
		for y := 0; y < 50; y++ {
			img.Set(x, y, color.White)
		}
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("encode png: %v", err)
	}

	text, err := runTesseractCLI(buf.Bytes())
	if err != nil {
		t.Fatalf("runTesseractCLI: %v", err)
	}
	_ = text
}

func TestOCRPDF_TesseractMissingHint(t *testing.T) {
	if tesseractAvailable() {
		t.Skip("tesseract is installed")
	}

	path := os.Getenv("OCR_TEST_SCANNED_PDF")
	if path == "" {
		t.Skip("OCR_TEST_SCANNED_PDF not set")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read test pdf: %v", err)
	}

	_, _, _, err = OCRPDF(data)
	if err == nil {
		t.Fatal("expected error when tesseract missing for scanned PDF")
	}
	if !containsAll(err.Error(), "Tesseract OCR kurulu değil", "brew install tesseract") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func containsAll(s string, parts ...string) bool {
	for _, p := range parts {
		if !bytes.Contains([]byte(s), []byte(p)) {
			return false
		}
	}
	return true
}
