package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestNormaliseArgs(t *testing.T) {
	var tests = []struct {
		input    []string
		expected []string
	}{
		{
			[]string{"-abc"},
			[]string{"-a", "-b", "-c"},
		},
		{
			[]string{"--hello", "--world"},
			[]string{"--hello", "--world"},
		},
		{
			[]string{"-a=b"},
			[]string{"-a", "b"},
		},
		{
			[]string{"-a", "b"},
			[]string{"-a", "b"},
		},
		{
			[]string{"--hello=world"},
			[]string{"--hello", "world"},
		},
		{
			[]string{"--hello", "world"},
			[]string{"--hello", "world"},
		},
		{
			[]string{"-a", "-bcd=e", "--hello=world", "--foo"},
			[]string{"-a", "-b", "-c", "-d", "e", "--hello", "world", "--foo"},
		},
		{
			[]string{"command", "-abc"},
			[]string{"command", "-a", "-b", "-c"},
		},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf(
				"should normalise %q correctly",
				strings.Join(tt.input, " "),
			),
			func(t *testing.T) {
				output := normaliseArgs(tt.input)

				if reflect.DeepEqual(output, tt.expected) == false {
					t.Errorf("Expected %#v, got %#v", tt.expected, output)
				}
			},
		)
	}
}
