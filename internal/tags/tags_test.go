package tags

import (
	"fmt"
	"reflect"
	"testing"
)

func TestExtractTags(t *testing.T) {
	t.Parallel()

	var testCases = []struct {
		input    string
		expected []string
	}{
		{"This is a #tag", []string{"tag"}},
		{"This is #another-tag", []string{"another-tag"}},
		{"This is #another/tag", []string{"another/tag"}},
		{"This is #another:tag", []string{"another:tag"}},
		{"This is a tag: #t", []string{"t"}},
		{"This is a #tag123", []string{"tag123"}},
		{"This is a tag: #123a", []string{"123a"}},
		{"This is a (#tag)", []string{"tag"}},
		{"This is a '#tag'", []string{"tag"}},
		{"This is a \"#tag\"", []string{"tag"}},
		{"This is \\#not-a-tag", []string(nil)},
		{"This is not#atag", []string(nil)},
		{"This is not a tag: #123", []string(nil)},
		{"This is not a tag: #:", []string(nil)},
		{"This is not a tag: #-", []string(nil)},
		{"This is not a tag: #/", []string(nil)},
		{"#tag", []string{"tag"}},
		{"#tag #another-tag #tag", []string{"tag", "another-tag", "tag"}},
		{"#tag #tag\nchicken #tag", []string{"tag", "tag", "tag"}},
	}

	for _, tt := range testCases {
		tt := tt

		t.Run(
			fmt.Sprintf("extracts tags from %q correctly", tt.input),
			func(t *testing.T) {
				t.Parallel()

				output := ExtractTags(tt.input)

				if reflect.DeepEqual(output, tt.expected) == false {
					t.Errorf("Expected %#v, got %#v", tt.expected, output)
				}
			},
		)
	}
}
