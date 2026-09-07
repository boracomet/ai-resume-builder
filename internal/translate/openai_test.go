package translate

import (
	"strings"
	"testing"
)

func TestParseNumberedTranslationResponseMultiline(t *testing.T) {
	reply := strings.Join([]string{
		"[1] Web Design and Coding",
		"[2] Anadolu University Open Education Faculty",
		"[3] 2022",
		"[4] Ongoing",
		"[5] I am continuing my education in the Web Design and Coding associate program.",
		"The program includes training in HTML, CSS, JavaScript, user experience (UX),",
		"web-based application development, and digital publishing.",
	}, "\n")

	results, err := parseNumberedTranslationResponse(reply, 5)
	if err != nil {
		t.Fatalf("parseNumberedTranslationResponse() error = %v", err)
	}

	if results[0] != "Web Design and Coding" {
		t.Fatalf("unexpected degree translation: %q", results[0])
	}
	if results[1] != "Anadolu University Open Education Faculty" {
		t.Fatalf("unexpected institution translation: %q", results[1])
	}
	if results[2] != "2022" {
		t.Fatalf("unexpected startDate translation: %q", results[2])
	}
	if results[4] == "" || !strings.Contains(results[4], "HTML, CSS, JavaScript") {
		t.Fatalf("unexpected description translation: %q", results[4])
	}
}

func TestTranslateEducationScalar(t *testing.T) {
	value, ok := translateEducationScalar("2022", "en")
	if !ok || value != "2022" {
		t.Fatalf("expected year passthrough, got %q ok=%v", value, ok)
	}

	value, ok = translateEducationScalar("Devam Ediyor", "en")
	if !ok || value != "Ongoing" {
		t.Fatalf("expected ongoing translation, got %q ok=%v", value, ok)
	}
}
