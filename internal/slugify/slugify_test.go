package slugify

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSlugify(t *testing.T) {
	var testCases = []struct {
		input    string
		expected string
	}{
		{"already-a-slug", "already-a-slug"},
		{"--already--a--slug--", "already-a-slug"},
		{"apple banana carrot", "apple-banana-carrot"},
		{"don't eat Brian’s \"cake\", please!", "dont-eat-brians-cake-please"},
		{"Crème Caramel with Äpfel", "crème-caramel-with-äpfel"},
	}

	for _, tt := range testCases {
		t.Run(
			fmt.Sprintf("slugifies %q correctly", tt.input),
			func(t *testing.T) {
				t.Parallel()

				output := Slugify(tt.input)

				if reflect.DeepEqual(output, tt.expected) == false {
					t.Errorf("Expected %#v, got %#v", tt.expected, output)
				}
			},
		)
	}
}
