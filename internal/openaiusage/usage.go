package openaiusage

import "strings"

// Usage holds token counts from an OpenAI chat completion response.
type Usage struct {
	PromptTokens     int `json:"promptTokens"`
	CompletionTokens int `json:"completionTokens"`
	TotalTokens      int `json:"totalTokens"`
}

func (u Usage) Add(other Usage) Usage {
	return Usage{
		PromptTokens:     u.PromptTokens + other.PromptTokens,
		CompletionTokens: u.CompletionTokens + other.CompletionTokens,
		TotalTokens:      u.TotalTokens + other.TotalTokens,
	}
}

type modelPricing struct {
	Input  float64 // USD per 1M tokens
	Output float64 // USD per 1M tokens
}

// Approximate pricing per 1M tokens (input, output).
var modelPricingTable = map[string]modelPricing{
	"gpt-4o":        {2.50, 10.00},
	"gpt-4o-mini":   {0.15, 0.60},
	"gpt-4.1":       {2.00, 8.00},
	"gpt-4.1-mini":  {0.40, 1.60},
	"gpt-4.1-nano":  {0.10, 0.40},
	"gpt-5":         {1.25, 10.00},
	"gpt-5-mini":    {0.25, 2.00},
	"gpt-5-nano":    {0.10, 0.40},
	"gpt-3.5-turbo": {0.50, 1.50},
	"o1":            {15.00, 60.00},
	"o1-mini":       {3.00, 12.00},
	"o3-mini":       {1.10, 4.40},
	"default":       {1.00, 3.00},
}

func lookupPricing(model string) modelPricing {
	model = strings.ToLower(strings.TrimSpace(model))
	if p, ok := modelPricingTable[model]; ok {
		return p
	}
	bestLen := 0
	var best modelPricing
	for prefix, p := range modelPricingTable {
		if prefix == "default" {
			continue
		}
		if strings.HasPrefix(model, prefix) && len(prefix) > bestLen {
			bestLen = len(prefix)
			best = p
		}
	}
	if bestLen > 0 {
		return best
	}
	return modelPricingTable["default"]
}

// CalculateCostUSD estimates request cost from token usage and model name.
func CalculateCostUSD(model string, u Usage) float64 {
	if u.PromptTokens == 0 && u.CompletionTokens == 0 {
		return 0
	}
	p := lookupPricing(model)
	inputCost := (float64(u.PromptTokens) / 1_000_000) * p.Input
	outputCost := (float64(u.CompletionTokens) / 1_000_000) * p.Output
	return inputCost + outputCost
}
