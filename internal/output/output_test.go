package output

import (
	"testing"
)

func TestMaskToken(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", "****"},
		{"short", "****"},
		{"12345678", "****"},
		{"123456789", "****23456789"},
		{"abcdefghijklmnop", "****ijklmnop"},
	}
	for _, tt := range tests {
		got := MaskToken(tt.input)
		if got != tt.expected {
			t.Errorf("MaskToken(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		input    string
		maxLen   int
		expected string
	}{
		{"hello", 10, "hello"},
		{"hello", 5, "hello"},
		{"hello world", 8, "hello..."},
		{"ab", 5, "ab"},
		{"abcdefghij", 6, "abc..."},
	}
	for _, tt := range tests {
		got := Truncate(tt.input, tt.maxLen)
		if got != tt.expected {
			t.Errorf("Truncate(%q, %d) = %q, want %q", tt.input, tt.maxLen, got, tt.expected)
		}
	}
}

func TestFormatList(t *testing.T) {
	tests := []struct {
		input    []string
		expected string
	}{
		{nil, ""},
		{[]string{}, ""},
		{[]string{"a"}, "a"},
		{[]string{"a", "b", "c"}, "a, b, c"},
	}
	for _, tt := range tests {
		got := FormatList(tt.input)
		if got != tt.expected {
			t.Errorf("FormatList(%v) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestNormalize(t *testing.T) {
	// Single map
	single := map[string]interface{}{"id": 1, "name": "test"}
	rows := normalize(single)
	if len(rows) != 1 {
		t.Errorf("expected 1 row, got %d", len(rows))
	}

	// Slice of maps
	multi := []map[string]interface{}{
		{"id": 1},
		{"id": 2},
	}
	rows = normalize(multi)
	if len(rows) != 2 {
		t.Errorf("expected 2 rows, got %d", len(rows))
	}

	// Slice of interface
	iface := []interface{}{
		map[string]interface{}{"id": 1},
		map[string]interface{}{"id": 2},
	}
	rows = normalize(iface)
	if len(rows) != 2 {
		t.Errorf("expected 2 rows from interface slice, got %d", len(rows))
	}

	// Nil - should not panic or stack overflow
	rows = normalize(nil)
	if len(rows) != 0 {
		t.Errorf("expected empty result for nil input, got %v", rows)
	}
}

func TestInferColumns(t *testing.T) {
	rows := []map[string]interface{}{
		{"name": "a", "id": 1},
		{"name": "b", "status": "ok"},
	}
	cols := inferColumns(rows)
	if len(cols) != 3 {
		t.Errorf("expected 3 columns, got %d: %v", len(cols), cols)
	}
	// Should be sorted
	if cols[0] != "id" || cols[1] != "name" || cols[2] != "status" {
		t.Errorf("unexpected column order: %v", cols)
	}
}
