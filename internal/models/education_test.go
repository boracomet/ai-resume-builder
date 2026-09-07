package models

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestEducationUnmarshalLegacySchool(t *testing.T) {
	raw := `{"degree":"Grafik Tasarım","school":"Nişantaşı Üniversitesi","startDate":"2018","endDate":"2020","description":"Açıklama"}`
	var edu Education
	if err := json.Unmarshal([]byte(raw), &edu); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if edu.Institution != "Nişantaşı Üniversitesi" {
		t.Fatalf("expected institution from legacy school, got %q", edu.Institution)
	}

	data, err := json.Marshal(edu)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}
	payload := string(data)
	if !strings.Contains(payload, `"institution":"Nişantaşı Üniversitesi"`) {
		t.Fatalf("unexpected marshaled education: %s", payload)
	}
	if strings.Contains(payload, `"school"`) {
		t.Fatalf("legacy school field should not be marshaled: %s", payload)
	}
}
