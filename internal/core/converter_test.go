package core

import "testing"

func TestConverter_ToCommand(t *testing.T) {
	c := NewConverter()

	tests := []struct {
		in  string
		out string
	}{
		{"list all directories", "ls -d */"},
		{"unknown", "# unknown command"},
	}

	for _, tt := range tests {
		got := c.ToCommand(tt.in)
		if got != tt.out {
			t.Errorf("expected %q got %q", tt.out, got)
		}
	}
}
