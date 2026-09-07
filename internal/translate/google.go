package translate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const translateEndpoint = "https://translation.googleapis.com/language/translate/v2"

type translateRequest struct {
	Q      []string `json:"q"`
	Source string   `json:"source"`
	Target string   `json:"target"`
	Format string   `json:"format"`
}

type translateResponse struct {
	Data struct {
		Translations []struct {
			TranslatedText string `json:"translatedText"`
		} `json:"translations"`
	} `json:"data"`
	Error *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"error"`
}

func TranslateTexts(texts []string, source, target, apiKey string) ([]string, error) {
	apiKey = strings.TrimSpace(apiKey)
	if apiKey == "" {
		return nil, fmt.Errorf("Google Translate API anahtarı gerekli (GOOGLE_TRANSLATE_API_KEY veya arayüzden girin)")
	}

	source = strings.TrimSpace(source)
	target = strings.TrimSpace(target)
	if source == "" || target == "" {
		return nil, fmt.Errorf("kaynak ve hedef dil gerekli")
	}
	if len(texts) == 0 {
		return nil, nil
	}

	body, err := json.Marshal(translateRequest{
		Q:      texts,
		Source: source,
		Target: target,
		Format: "text",
	})
	if err != nil {
		return nil, err
	}

	endpoint := translateEndpoint + "?" + url.Values{"key": {apiKey}}.Encode()
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("çeviri isteği başarısız: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var parsed translateResponse
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return nil, fmt.Errorf("çeviri yanıtı okunamadı: %w", err)
	}

	if parsed.Error != nil {
		return nil, fmt.Errorf("Google Translate hatası: %s", parsed.Error.Message)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Google Translate HTTP %d", resp.StatusCode)
	}

	translations := parsed.Data.Translations
	if len(translations) != len(texts) {
		return nil, fmt.Errorf("çeviri sonuç sayısı uyuşmuyor: beklenen %d, alınan %d", len(texts), len(translations))
	}

	results := make([]string, len(translations))
	for i, t := range translations {
		results[i] = t.TranslatedText
	}
	return results, nil
}
