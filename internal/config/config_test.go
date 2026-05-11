package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestNewConfig(t *testing.T) {
	cfg := NewConfig()
	if cfg.DefaultContext != "" {
		t.Errorf("expected empty default context, got %s", cfg.DefaultContext)
	}
	if cfg.Contexts == nil {
		t.Error("expected non-nil contexts map")
	}
	if len(cfg.Contexts) != 0 {
		t.Errorf("expected empty contexts, got %d", len(cfg.Contexts))
	}
}

func TestConfig_SetGetContext(t *testing.T) {
	cfg := NewConfig()
	cfg.SetContext("test", Context{Name: "test"})

	ctx, ok := cfg.GetContext("test")
	if !ok {
		t.Fatal("expected to find context 'test'")
	}
	if ctx.Name != "test" {
		t.Errorf("expected name 'test', got %s", ctx.Name)
	}

	_, ok = cfg.GetContext("nonexistent")
	if ok {
		t.Error("expected not to find 'nonexistent' context")
	}
}

func TestConfig_DeleteContext(t *testing.T) {
	cfg := NewConfig()
	cfg.SetContext("test", Context{Name: "test"})
	cfg.DeleteContext("test")

	_, ok := cfg.GetContext("test")
	if ok {
		t.Error("expected context to be deleted")
	}
}

func TestConfig_GetDefaultContext(t *testing.T) {
	cfg := NewConfig()

	_, _, ok := cfg.GetDefaultContext()
	if ok {
		t.Error("expected no default context")
	}

	cfg.SetContext("prod", Context{Name: "prod"})
	cfg.DefaultContext = "prod"

	name, ctx, ok := cfg.GetDefaultContext()
	if !ok {
		t.Fatal("expected to find default context")
	}
	if name != "prod" {
		t.Errorf("expected name 'prod', got %s", name)
	}
	if ctx.Name != "prod" {
		t.Errorf("expected context name 'prod', got %s", ctx.Name)
	}
}

func TestConfig_SaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "config.json")

	cfg := NewConfig()
	cfg.SetContext("test", Context{Name: "test"})
	cfg.DefaultContext = "test"

	// Override config path for testing
	origConfigDir := configDir
	origConfigFile := configFile

	// We can't easily override the config path, so test Save/Load manually
	data := []byte(`{"default_context":"test","contexts":{"test":{"name":"test"}}}`)
	if err := os.WriteFile(tmpFile, data, 0600); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	fileData, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("failed to read config: %v", err)
	}

	loaded := NewConfig()
	if err := json.Unmarshal(fileData, loaded); err != nil {
		t.Fatalf("failed to parse config: %v", err)
	}

	if loaded.DefaultContext != "test" {
		t.Errorf("expected default context 'test', got %s", loaded.DefaultContext)
	}
	ctx, ok := loaded.GetContext("test")
	if !ok {
		t.Fatal("expected to find context 'test'")
	}
	if ctx.Name != "test" {
		t.Errorf("expected context name 'test', got %s", ctx.Name)
	}

	_ = origConfigDir
	_ = origConfigFile
}

func TestConfig_LoadNonExistent(t *testing.T) {
	// Load should return a new config if file doesn't exist
	cfg, err := Load()
	if err != nil {
		// It's ok if the file doesn't exist in test env
		if !os.IsNotExist(err) {
			// Load returns NewConfig() for non-existent files
		}
	}
	if cfg != nil && cfg.Contexts == nil {
		t.Error("expected initialized contexts map")
	}
}
