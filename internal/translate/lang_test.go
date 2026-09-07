package translate

import "testing"

func TestValidLangCode(t *testing.T) {
	tests := []struct {
		code string
		ok   bool
	}{
		{"tr", true},
		{"en", true},
		{"de", true},
		{"DE", true},
		{" fr ", true},
		{"eng", false},
		{"", false},
		{"1a", false},
	}
	for _, tc := range tests {
		if ValidLangCode(tc.code) != tc.ok {
			t.Errorf("ValidLangCode(%q) = %v, want %v", tc.code, !tc.ok, tc.ok)
		}
	}
}

func TestStorageSupported(t *testing.T) {
	if !StorageSupported("tr") || !StorageSupported("en") {
		t.Fatal("tr and en should be storage supported")
	}
	if !StorageSupported("de") || !StorageSupported("tl") {
		t.Fatal("valid ISO 639-1 codes should be storage supported")
	}
	if StorageSupported("eng") || StorageSupported("") {
		t.Fatal("invalid codes should not be storage supported")
	}
}

func TestResolveSourceLang(t *testing.T) {
	if got := ResolveSourceLang("de", "tr", "en"); got != "de" {
		t.Fatalf("explicit source: got %q", got)
	}
	if got := ResolveSourceLang("", "en", "tr"); got != "en" {
		t.Fatalf("editing lang: got %q", got)
	}
	if got := ResolveSourceLang("", "", "en"); got != "tr" {
		t.Fatalf("fallback for en target: got %q", got)
	}
}
