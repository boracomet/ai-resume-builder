package translate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/boracomet/ai-resume-builder/internal/openaiusage"
)

const openAIBaseURL = "https://api.openai.com/v1"

var numberedTranslationLineRe = regexp.MustCompile(`^\[(\d+)\]\s*(.*)$`)

func TranslateTextsWithOpenAI(texts []string, source, target, apiKey, model string) ([]string, openaiusage.Usage, float64, error) {
	return translateWithOpenAI(apiKey, model, texts, source, target)
}

func translateWithOpenAI(apiKey, model string, texts []string, source, target string) ([]string, openaiusage.Usage, float64, error) {
	if len(texts) == 0 {
		return nil, openaiusage.Usage{}, 0, nil
	}

	apiKey = strings.TrimSpace(apiKey)
	model = strings.TrimSpace(model)
	if apiKey == "" {
		return nil, openaiusage.Usage{}, 0, fmt.Errorf("OpenAI API anahtarı gerekli (OPENAI_API_KEY veya arayüzden girin)")
	}
	if model == "" {
		return nil, openaiusage.Usage{}, 0, fmt.Errorf("OpenAI model seçimi gerekli")
	}

	source = strings.TrimSpace(source)
	target = strings.TrimSpace(target)
	if source == "" || target == "" {
		return nil, openaiusage.Usage{}, 0, fmt.Errorf("kaynak ve hedef dil gerekli")
	}

	var numbered strings.Builder
	for i, text := range texts {
		if i > 0 {
			numbered.WriteString("\n")
		}
		numbered.WriteString(fmt.Sprintf("[%d] %s", i+1, text))
	}

	prompt := fmt.Sprintf(
		"Translate each numbered line from %s to %s. Keep the [N] prefix and return one translated line per original line. Preserve formatting and proper nouns when appropriate.\n\n%s",
		source, target, numbered.String(),
	)

	reply, usage, cost, err := openAIChat(apiKey, model, []openAIMessage{
		{Role: "system", Content: "You are a professional translator. Return only the translated numbered lines."},
		{Role: "user", Content: prompt},
	})
	if err != nil {
		return nil, openaiusage.Usage{}, 0, err
	}

	results, err := parseNumberedTranslationResponse(reply, len(texts))
	if err != nil {
		return nil, openaiusage.Usage{}, 0, err
	}
	return results, usage, cost, nil
}

func parseNumberedTranslationResponse(reply string, expectedCount int) ([]string, error) {
	if expectedCount == 0 {
		return nil, nil
	}

	results := make([]string, expectedCount)
	filled := make([]bool, expectedCount)
	currentIndex := -1
	var currentBuilder strings.Builder

	flush := func() {
		if currentIndex < 0 || currentIndex >= expectedCount {
			return
		}
		results[currentIndex] = strings.TrimSpace(currentBuilder.String())
		filled[currentIndex] = true
		currentBuilder.Reset()
	}

	for _, line := range strings.Split(reply, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" && currentIndex < 0 {
			continue
		}

		if matches := numberedTranslationLineRe.FindStringSubmatch(trimmed); matches != nil {
			flush()
			index, err := strconv.Atoi(matches[1])
			if err != nil || index < 1 || index > expectedCount {
				return nil, fmt.Errorf("çeviri yanıtında geçersiz satır numarası: %s", matches[1])
			}
			currentIndex = index - 1
			if rest := strings.TrimSpace(matches[2]); rest != "" {
				currentBuilder.WriteString(rest)
			}
			continue
		}

		if currentIndex < 0 {
			continue
		}
		if currentBuilder.Len() > 0 {
			currentBuilder.WriteString(" ")
		}
		currentBuilder.WriteString(trimmed)
	}
	flush()

	for i := range filled {
		if !filled[i] {
			return nil, fmt.Errorf("çeviri yanıtında eksik satır: [%d]", i+1)
		}
	}
	return results, nil
}

type openAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIChatRequest struct {
	Model    string          `json:"model"`
	Messages []openAIMessage `json:"messages"`
}

type openAIChatResponse struct {
	Choices []struct {
		Message openAIMessage `json:"message"`
	} `json:"choices"`
	Usage *struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

func openAIChat(apiKey, model string, messages []openAIMessage) (string, openaiusage.Usage, float64, error) {
	payload := openAIChatRequest{
		Model:    model,
		Messages: messages,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", openaiusage.Usage{}, 0, err
	}

	req, err := http.NewRequest(http.MethodPost, openAIBaseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", openaiusage.Usage{}, 0, err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", openaiusage.Usage{}, 0, fmt.Errorf("OpenAI isteği başarısız: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", openaiusage.Usage{}, 0, err
	}

	var parsed openAIChatResponse
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return "", openaiusage.Usage{}, 0, fmt.Errorf("OpenAI yanıtı okunamadı: %w", err)
	}
	if parsed.Error != nil {
		return "", openaiusage.Usage{}, 0, fmt.Errorf("OpenAI hatası: %s", parsed.Error.Message)
	}
	if resp.StatusCode != http.StatusOK {
		return "", openaiusage.Usage{}, 0, fmt.Errorf("OpenAI HTTP %d", resp.StatusCode)
	}
	if len(parsed.Choices) == 0 {
		return "", openaiusage.Usage{}, 0, fmt.Errorf("OpenAI boş yanıt döndü")
	}

	content := strings.TrimSpace(parsed.Choices[0].Message.Content)
	if content == "" {
		return "", openaiusage.Usage{}, 0, fmt.Errorf("OpenAI boş içerik döndü")
	}

	usage := openaiusage.Usage{}
	if parsed.Usage != nil {
		usage = openaiusage.Usage{
			PromptTokens:     parsed.Usage.PromptTokens,
			CompletionTokens: parsed.Usage.CompletionTokens,
			TotalTokens:      parsed.Usage.TotalTokens,
		}
	}
	cost := openaiusage.CalculateCostUSD(model, usage)
	return content, usage, cost, nil
}
