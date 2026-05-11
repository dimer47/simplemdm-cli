package keychain

import (
	"testing"
)

func TestKeyFor(t *testing.T) {
	tests := []struct {
		context  string
		expected string
	}{
		{"default", "simplemdm-cli:default"},
		{"prod", "simplemdm-cli:prod"},
		{"my-context", "simplemdm-cli:my-context"},
	}
	for _, tt := range tests {
		got := keyFor(tt.context)
		if got != tt.expected {
			t.Errorf("keyFor(%q) = %q, want %q", tt.context, got, tt.expected)
		}
	}
}

func TestKeychain_Integration(t *testing.T) {
	if !IsAvailable() {
		t.Skip("system keychain not available in this environment")
	}

	ctx := "__simplemdm_cli_test__"
	testKey := "test-api-key-12345"

	// Set
	if err := Set(ctx, testKey); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Get
	got, err := Get(ctx)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if got != testKey {
		t.Errorf("Get(%q) = %q, want %q", ctx, got, testKey)
	}

	// Delete
	if err := Delete(ctx); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Get after delete should fail
	_, err = Get(ctx)
	if err == nil {
		t.Error("expected error after delete, got nil")
	}
}
