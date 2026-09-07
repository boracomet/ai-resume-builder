package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/boracomet/ai-resume-builder/internal/openaiusage"
)

const openAIBaseURL = "https://api.openai.com/v1"

type Message struct {
	Role       string     `json:"role"`
	Content    string     `json:"content,omitempty"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
	ToolCallID string     `json:"tool_call_id,omitempty"`
}

type ToolCall struct {
	ID       string       `json:"id"`
	Type     string       `json:"type"`
	Function ToolFunction `json:"function"`
}

type ToolFunction struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type ToolDefinition struct {
	Type     string             `json:"type"`
	Function ToolFunctionSchema `json:"function"`
}

type ToolFunctionSchema struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

type modelListResponse struct {
	Data []struct {
		ID string `json:"id"`
	} `json:"data"`
	Error *openAIError `json:"error"`
}

type chatRequest struct {
	Model          string           `json:"model"`
	Messages       []Message        `json:"messages"`
	Temperature    *float64         `json:"temperature,omitempty"`
	ResponseFormat *responseFormat  `json:"response_format,omitempty"`
	Tools          []ToolDefinition `json:"tools,omitempty"`
	ToolChoice     interface{}      `json:"tool_choice,omitempty"`
}

func supportsCustomTemperature(model string) bool {
	model = strings.ToLower(strings.TrimSpace(model))
	restrictedPrefixes := []string{"o1", "o3", "o4", "gpt-5"}
	for _, prefix := range restrictedPrefixes {
		if strings.HasPrefix(model, prefix) {
			return false
		}
	}
	return true
}

type responseFormat struct {
	Type string `json:"type"`
}

type chatResponse struct {
	Choices []struct {
		Message      Message `json:"message"`
		FinishReason string  `json:"finish_reason"`
	} `json:"choices"`
	Usage *struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Error *openAIError `json:"error"`
}

type ChatCompletionResult struct {
	Message      Message
	FinishReason string
	Usage        openaiusage.Usage
	CostUSD      float64
}

type ToolChatResult struct {
	Reply   string
	Usage   openaiusage.Usage
	CostUSD float64
}

type TestConnectionResult struct {
	Usage   openaiusage.Usage
	CostUSD float64
}

type openAIError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    string `json:"code"`
}

func ResolveAPIKey(fromBody, fromQuery, fromHeader string) string {
	for _, candidate := range []string{fromBody, fromQuery, fromHeader} {
		if key := strings.TrimSpace(candidate); key != "" {
			return key
		}
	}
	return strings.TrimSpace(os.Getenv("OPENAI_API_KEY"))
}

func isChatCapableModel(id string) bool {
	id = strings.ToLower(strings.TrimSpace(id))
	if id == "" {
		return false
	}
	prefixes := []string{"gpt-4", "gpt-3.5", "gpt-5", "o1", "o3", "o4", "chatgpt"}
	for _, prefix := range prefixes {
		if strings.HasPrefix(id, prefix) {
			return true
		}
	}
	return false
}

func ListModels(apiKey string) ([]string, error) {
	apiKey = strings.TrimSpace(apiKey)
	if apiKey == "" {
		return nil, fmt.Errorf("OpenAI API anahtarı gerekli (OPENAI_API_KEY veya arayüzden girin)")
	}

	req, err := http.NewRequest(http.MethodGet, openAIBaseURL+"/models", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("OpenAI modelleri alınamadı: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var parsed modelListResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, fmt.Errorf("OpenAI yanıtı okunamadı: %w", err)
	}
	if parsed.Error != nil {
		return nil, fmt.Errorf("OpenAI hatası: %s", parsed.Error.Message)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OpenAI HTTP %d", resp.StatusCode)
	}

	models := make([]string, 0)
	for _, item := range parsed.Data {
		if isChatCapableModel(item.ID) {
			models = append(models, item.ID)
		}
	}
	sort.Strings(models)
	return models, nil
}

func TestConnection(apiKey, model string) (TestConnectionResult, error) {
	apiKey = strings.TrimSpace(apiKey)
	model = strings.TrimSpace(model)
	if apiKey == "" {
		return TestConnectionResult{}, fmt.Errorf("OpenAI API anahtarı gerekli (OPENAI_API_KEY veya arayüzden girin)")
	}
	if model == "" {
		return TestConnectionResult{}, fmt.Errorf("model seçimi gerekli")
	}

	result, err := ChatCompletion(apiKey, model, []Message{
		{Role: "user", Content: "ping"},
	}, nil, false)
	if err != nil {
		return TestConnectionResult{}, err
	}
	return TestConnectionResult{
		Usage:   result.Usage,
		CostUSD: result.CostUSD,
	}, nil
}

func Chat(apiKey, model string, messages []Message, jsonMode bool) (string, error) {
	result, err := ChatCompletion(apiKey, model, messages, nil, jsonMode)
	if err != nil {
		return "", err
	}
	content := strings.TrimSpace(result.Message.Content)
	if content == "" && len(result.Message.ToolCalls) == 0 {
		return "", fmt.Errorf("OpenAI boş içerik döndü")
	}
	return content, nil
}

func ChatCompletion(apiKey, model string, messages []Message, tools []ToolDefinition, jsonMode bool) (ChatCompletionResult, error) {
	apiKey = strings.TrimSpace(apiKey)
	model = strings.TrimSpace(model)
	if apiKey == "" {
		return ChatCompletionResult{}, fmt.Errorf("OpenAI API anahtarı gerekli (OPENAI_API_KEY veya arayüzden girin)")
	}
	if model == "" {
		return ChatCompletionResult{}, fmt.Errorf("model seçimi gerekli")
	}
	if len(messages) == 0 {
		return ChatCompletionResult{}, fmt.Errorf("mesaj listesi boş olamaz")
	}

	payload := chatRequest{
		Model:    model,
		Messages: messages,
	}
	if supportsCustomTemperature(model) {
		temp := 0.7
		payload.Temperature = &temp
	}
	if len(tools) > 0 {
		payload.Tools = tools
		payload.ToolChoice = "auto"
	} else if jsonMode {
		payload.ResponseFormat = &responseFormat{Type: "json_object"}
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return ChatCompletionResult{}, err
	}

	req, err := http.NewRequest(http.MethodPost, openAIBaseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return ChatCompletionResult{}, err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return ChatCompletionResult{}, fmt.Errorf("OpenAI isteği başarısız: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return ChatCompletionResult{}, err
	}

	var parsed chatResponse
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return ChatCompletionResult{}, fmt.Errorf("OpenAI yanıtı okunamadı: %w", err)
	}
	if parsed.Error != nil {
		return ChatCompletionResult{}, fmt.Errorf("OpenAI hatası: %s", parsed.Error.Message)
	}
	if resp.StatusCode != http.StatusOK {
		return ChatCompletionResult{}, fmt.Errorf("OpenAI HTTP %d", resp.StatusCode)
	}
	if len(parsed.Choices) == 0 {
		return ChatCompletionResult{}, fmt.Errorf("OpenAI boş yanıt döndü")
	}

	choice := parsed.Choices[0]
	if strings.TrimSpace(choice.Message.Content) == "" && len(choice.Message.ToolCalls) == 0 {
		return ChatCompletionResult{}, fmt.Errorf("OpenAI boş içerik döndü")
	}

	usage := openaiusage.Usage{}
	if parsed.Usage != nil {
		usage = openaiusage.Usage{
			PromptTokens:     parsed.Usage.PromptTokens,
			CompletionTokens: parsed.Usage.CompletionTokens,
			TotalTokens:      parsed.Usage.TotalTokens,
		}
	}

	return ChatCompletionResult{
		Message:      choice.Message,
		FinishReason: choice.FinishReason,
		Usage:        usage,
		CostUSD:      openaiusage.CalculateCostUSD(model, usage),
	}, nil
}

func RunToolChatLoop(
	apiKey, model string,
	messages []Message,
	executor *ToolExecutor,
	ctx *ToolContext,
	state *ToolExecutionState,
) (ToolChatResult, error) {
	tools := ToolDefinitions()
	totalUsage := openaiusage.Usage{}
	totalCost := 0.0

	for i := 0; i < maxToolIterations; i++ {
		result, err := ChatCompletion(apiKey, model, messages, tools, false)
		if err != nil {
			return ToolChatResult{}, err
		}
		totalUsage = totalUsage.Add(result.Usage)
		totalCost += result.CostUSD

		if len(result.Message.ToolCalls) == 0 {
			return ToolChatResult{
				Reply:   strings.TrimSpace(result.Message.Content),
				Usage:   totalUsage,
				CostUSD: totalCost,
			}, nil
		}

		messages = append(messages, result.Message)

		for _, toolCall := range result.Message.ToolCalls {
			var args map[string]interface{}
			if strings.TrimSpace(toolCall.Function.Arguments) != "" {
				if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); err != nil {
					args = map[string]interface{}{}
				}
			}
			if args == nil {
				args = map[string]interface{}{}
			}

			toolResult, execErr := executor.Execute(toolCall.Function.Name, args, ctx, state)
			if execErr != nil && toolResult == "" {
				toolResult = fmt.Sprintf(`{"error":%q}`, execErr.Error())
			}
			totalUsage = totalUsage.Add(state.LastUsage)
			totalCost += state.LastCostUSD
			state.LastUsage = openaiusage.Usage{}
			state.LastCostUSD = 0

			messages = append(messages, Message{
				Role:       "tool",
				ToolCallID: toolCall.ID,
				Content:    toolResult,
			})
		}
	}

	return ToolChatResult{}, fmt.Errorf("araç çağrısı döngüsü limiti aşıldı (%d)", maxToolIterations)
}
