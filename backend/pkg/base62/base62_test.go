package base62

import (
	"testing"
)

func TestEncode(t *testing.T) {
	tests := []struct {
		name     string
		input    int64
		expected string
	}{
		{"Zero", 0, "0"},
		{"Single digit", 1, "1"},
		{"Base boundary", 61, "Z"},
		{"Multi-character", 62, "10"},
		{"Arbitrary number", 12345, "3d7"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Encode(tt.input)
			if got != tt.expected {
				t.Errorf("Encode(%d) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestGenerateRandom(t *testing.T) {
	length := 8
	code1, err := GenerateRandom(length)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(code1) != length {
		t.Errorf("Expected length %d, got %d", length, len(code1))
	}

	code2, err := GenerateRandom(length)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if code1 == code2 {
		t.Errorf("Expected distinct random strings, got identical %s", code1)
	}
}
