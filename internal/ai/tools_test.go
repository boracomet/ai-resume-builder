package ai

import (
	"path/filepath"
	"testing"

	"github.com/boracomet/ai-resume-builder/internal/db"
)

func setupTestRepo(t *testing.T) *db.CVRepository {
	t.Helper()
	conn, err := db.Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { conn.Close() })
	return db.NewCVRepository(conn)
}

func TestToolExecutorCreateProfile(t *testing.T) {
	repo := setupTestRepo(t)
	executor := NewToolExecutor(repo)
	ctx := &ToolContext{Language: "tr"}
	state := &ToolExecutionState{}

	resultJSON, err := executor.Execute("create_profile", map[string]interface{}{
		"name":     "Test Engineer CV",
		"language": "tr",
	}, ctx, state)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !state.ProfileUpdated {
		t.Fatal("expected profileUpdated=true")
	}
	if state.Profile == nil || state.Profile.Name != "Test Engineer CV" {
		t.Fatalf("unexpected profile: %+v", state.Profile)
	}
	if ctx.ProfileID == nil || *ctx.ProfileID != state.Profile.ID {
		t.Fatal("context profileId not updated")
	}
	if len(state.Actions) != 1 || state.Actions[0].Tool != "create_profile" {
		t.Fatalf("unexpected actions: %+v", state.Actions)
	}
	if resultJSON == "" {
		t.Fatal("expected non-empty tool result")
	}
}

func TestToolExecutorSwitchLanguage(t *testing.T) {
	repo := setupTestRepo(t)
	executor := NewToolExecutor(repo)
	ctx := &ToolContext{Language: "tr"}
	state := &ToolExecutionState{}

	_, err := executor.Execute("create_profile", map[string]interface{}{
		"name": "Lang Test",
	}, ctx, state)
	if err != nil {
		t.Fatal(err)
	}

	state = &ToolExecutionState{}
	_, err = executor.Execute("switch_language", map[string]interface{}{
		"language": "en",
	}, ctx, state)
	if err != nil {
		t.Fatalf("switch_language error = %v", err)
	}
	if state.Language != "en" {
		t.Fatalf("expected language en, got %s", state.Language)
	}
	if state.Profile == nil || state.Profile.Language != "en" {
		t.Fatalf("profile language not updated: %+v", state.Profile)
	}
}

func TestToolDefinitionsCount(t *testing.T) {
	tools := ToolDefinitions()
	if len(tools) != 9 {
		t.Fatalf("expected 9 tools, got %d", len(tools))
	}
	names := map[string]bool{}
	for _, tool := range tools {
		names[tool.Function.Name] = true
	}
	required := []string{
		"create_profile",
		"duplicate_profile",
		"update_profile_content",
		"translate_profile",
		"switch_language",
		"get_current_profile",
		"apply_to_profile",
		"request_profile_selection",
		"copy_from_profile",
	}
	for _, name := range required {
		if !names[name] {
			t.Fatalf("missing tool definition: %s", name)
		}
	}
}
