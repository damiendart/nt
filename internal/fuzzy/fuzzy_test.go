package fuzzy

import (
	"fmt"
	"testing"
)

func TestIsFuzzyMatch(t *testing.T) {
	var testCases = []struct {
		needle   string
		haystack string
		expected bool
	}{
		{"hell", "hello", true},
		{"ho", "hello", true},
		{"al", "äpfel", true},
		{"HELL", "hello", true},
		{"hlo", "HELLO", true},

		{"hill", "hello", false},
		{"goodbye", "hello", false},
	}

	for _, tt := range testCases {
		t.Run(
			fmt.Sprintf("%q matches %q correctly", tt.needle, tt.haystack),
			func(t *testing.T) {
				t.Parallel()

				output := IsFuzzyMatch(tt.needle, tt.haystack)

				if output != tt.expected {
					t.Errorf("Expected %#v, got %#v", tt.expected, output)
				}
			},
		)
	}
}
