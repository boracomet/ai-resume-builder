package ocr

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const tesseractInstallHint = "Tesseract OCR kurulu değil; taranmış PDF'ler için: brew install tesseract tesseract-lang"

func tesseractAvailable() bool {
	_, err := exec.LookPath("tesseract")
	return err == nil
}

func runTesseractCLI(pngData []byte) (string, error) {
	tmpDir, err := os.MkdirTemp("", "bora-cv-ocr-*")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tmpDir)

	imgPath := filepath.Join(tmpDir, "page.png")
	if err := os.WriteFile(imgPath, pngData, 0o600); err != nil {
		return "", err
	}

	cmd := exec.Command("tesseract", imgPath, "stdout", "-l", "tur+eng", "--psm", "3")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg != "" {
			return "", fmt.Errorf("%w: %s", err, msg)
		}
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}
