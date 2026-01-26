package ui

import (
	"strings"
	"testing"
)

func TestCheckmark(t *testing.T) {
	result := Checkmark()
	if result == "" {
		t.Error("Checkmark() returned empty string")
	}
	// Should contain the checkmark character
	if !strings.Contains(result, "\u2713") {
		t.Errorf("Checkmark() = %q, expected to contain checkmark character", result)
	}
}

func TestCross(t *testing.T) {
	result := Cross()
	if result == "" {
		t.Error("Cross() returned empty string")
	}
	// Should contain the X character
	if !strings.Contains(result, "\u2717") {
		t.Errorf("Cross() = %q, expected to contain X character", result)
	}
}

func TestMaskAPIKey(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantLast string
	}{
		{
			name:     "empty key",
			input:    "",
			wantLast: "",
		},
		{
			name:     "short key",
			input:    "abc",
			wantLast: "",
		},
		{
			name:     "normal key",
			input:    "sk-1234567890abcd",
			wantLast: "abcd",
		},
		{
			name:     "exactly 4 chars",
			input:    "abcd",
			wantLast: "",
		},
		{
			name:     "5 chars",
			input:    "abcde",
			wantLast: "bcde",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MaskAPIKey(tt.input)
			if tt.wantLast != "" {
				if !strings.HasSuffix(result, tt.wantLast) {
					t.Errorf("MaskAPIKey(%q) = %q, expected to end with %q", tt.input, result, tt.wantLast)
				}
			}
			// Should contain bullet characters if key is non-empty
			if tt.input != "" && len(tt.input) > 4 {
				if !strings.Contains(result, "\u2022") {
					t.Errorf("MaskAPIKey(%q) = %q, expected to contain bullet characters", tt.input, result)
				}
			}
		})
	}
}
