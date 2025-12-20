package fuzzy

import (
	"fmt"
	"testing"
)

func TestMatchScore(t *testing.T) {
	var testCases = []struct {
		s        string
		needle   string
		expected int
	}{
		{"hello", "hell", 4},
		{"hello", "ho", 5},
		{"hello", "eo", 5},
		{"hello", "HELL", 4},
		{"HELLO", "hlo", 5},
		{"hello", "o", 5},
		{"äpfel", "al", 5},

		{"hello", "", -1},
		{"", "hello", -1},
		{"hello", "hill", -1},
		{"hello", "goodbye", -1},
		{"greet", "greetings", -1},
	}

	for _, tt := range testCases {
		t.Run(
			fmt.Sprintf("fuzzy matching %q against %q returns the correct score", tt.needle, tt.s),
			func(t *testing.T) {
				t.Parallel()

				output := MatchScore(tt.s, tt.needle)

				if output != tt.expected {
					t.Errorf("Expected %#v, got %#v", tt.expected, output)
				}
			},
		)
	}
}
