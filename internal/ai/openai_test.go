package ai

import (
	"encoding/json"
	"testing"
)

func TestSupportsCustomTemperature(t *testing.T) {
	cases := []struct {
		model    string
		supports bool
	}{
		{"gpt-4o-mini", true},
		{"gpt-4o", true},
		{"gpt-3.5-turbo", true},
		{"gpt-5-nano", false},
		{"gpt-5-mini", false},
		{"gpt-5", false},
		{"o1-preview", false},
		{"o1-mini", false},
		{"o3-mini", false},
		{"o4-mini", false},
	}
	for _, tc := range cases {
		got := supportsCustomTemperature(tc.model)
		if got != tc.supports {
			t.Errorf("supportsCustomTemperature(%q) = %v, want %v", tc.model, got, tc.supports)
		}
	}
}

func TestChatRequestOmitsTemperatureForRestrictedModels(t *testing.T) {
	payload := chatRequest{
		Model:    "gpt-5-nano",
		Messages: []Message{{Role: "user", Content: "ping"}},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != `{"model":"gpt-5-nano","messages":[{"role":"user","content":"ping"}]}` {
		t.Errorf("unexpected payload: %s", body)
	}

	temp := 0.7
	payload.Temperature = &temp
	body, err = json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != `{"model":"gpt-5-nano","messages":[{"role":"user","content":"ping"}],"temperature":0.7}` {
		t.Errorf("unexpected payload with temperature: %s", body)
	}
}
