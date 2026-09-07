package openaiusage

import (
	"math"
	"testing"
)

func TestCalculateCostUSD(t *testing.T) {
	u := Usage{PromptTokens: 1_000_000, CompletionTokens: 0}
	cost := CalculateCostUSD("gpt-4o-mini", u)
	if math.Abs(cost-0.15) > 1e-9 {
		t.Errorf("gpt-4o-mini input cost = %v, want 0.15", cost)
	}

	u = Usage{PromptTokens: 0, CompletionTokens: 1_000_000}
	cost = CalculateCostUSD("gpt-4o-mini", u)
	if math.Abs(cost-0.60) > 1e-9 {
		t.Errorf("gpt-4o-mini output cost = %v, want 0.60", cost)
	}

	u = Usage{PromptTokens: 100, CompletionTokens: 50}
	cost = CalculateCostUSD("gpt-4o", u)
	expected := (100.0/1e6)*2.50 + (50.0/1e6)*10.00
	if math.Abs(cost-expected) > 1e-12 {
		t.Errorf("gpt-4o mixed cost = %v, want %v", cost, expected)
	}
}

func TestUsageAdd(t *testing.T) {
	a := Usage{PromptTokens: 10, CompletionTokens: 5, TotalTokens: 15}
	b := Usage{PromptTokens: 20, CompletionTokens: 10, TotalTokens: 30}
	sum := a.Add(b)
	if sum.PromptTokens != 30 || sum.CompletionTokens != 15 || sum.TotalTokens != 45 {
		t.Errorf("unexpected sum: %+v", sum)
	}
}
