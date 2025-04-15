// Copyright (C) Damien Dart, <damiendart@pobox.com>.
// This file is distributed under the MIT licence. For more information,
// please refer to the accompanying "LICENCE" file.

package cli

import (
	"fmt"
	"reflect"
	"testing"
)

func TestParseArgsWithInvalidInput(t *testing.T) {
	t.Parallel()

	var testCases = []struct {
		spec          Spec
		input         []string
		expectedError error
	}{
		{
			Spec{},
			[]string{"--alpha"},
			UnknownOptionError("alpha"),
		},
		{
			Spec{"alpha": ValueRequired},
			[]string{"--alpha"},
			MissingOptionValueError("alpha"),
		},
		{
			Spec{"alpha": NoValueAccepted},
			[]string{"--alpha=B"},
			UnexpectedOptionValueError("alpha"),
		},
		{
			Spec{"alpha": ValueRequired},
			[]string{"--alpha", "--", "BETA"},
			MissingOptionValueError("alpha"),
		},
	}

	for _, tt := range testCases {
		t.Run(
			fmt.Sprintf("handles %q correctly", tt.input),
			func(t *testing.T) {
				t.Parallel()

				_, _, err := ParseArgs(tt.input, tt.spec)
				if err == nil {
					t.Errorf("expected %#v, got %#v", tt.expectedError.Error(), err)
				} else if err.Error() != tt.expectedError.Error() {
					t.Errorf("expected %#v, got %#v", tt.expectedError.Error(), err.Error())
				}
			},
		)
	}

}

func TestParseArgsWithValidInputs(t *testing.T) {
	t.Parallel()

	var testCases = []struct {
		spec                       Spec
		input                      []string
		expectedOptions            OptionMap
		expectedRemainingArguments []string
	}{
		{
			Spec{},
			[]string{},
			OptionMap{},
			[]string{},
		},
		{
			Spec{"a": NoValueAccepted, "b": NoValueAccepted, "c": ValueRequired, "e": ValueRequired, "h": ValueOptional, "j": ValueOptional},
			[]string{"-abcALPHA BETA", "-e=FG", "-h-I", "-j", "K"},
			OptionMap{"a": "", "b": "", "c": "ALPHA BETA", "e": "=FG", "h": "-I", "j": ""},
			[]string{"K"},
		},
		{
			Spec{"alpha": ValueRequired, "gamma": ValueRequired, "epsilon": NoValueAccepted, "theta": ValueOptional, "kappa": ValueOptional},
			[]string{"--alpha", "BETA", "--gamma=DELTA", "--epsilon", "--theta=--iota", "--kappa", "lambda"},
			OptionMap{"alpha": "BETA", "gamma": "DELTA", "epsilon": "", "theta": "--iota", "kappa": ""},
			[]string{"lambda"},
		},
		{
			Spec{"a": NoValueAccepted, "b": NoValueAccepted, "c": ValueRequired, "e": ValueRequired, "alpha": ValueOptional, "beta": NoValueAccepted},
			[]string{"-abc", "d", "--alpha", "-eFG", "--beta", "g"},
			OptionMap{"a": "", "b": "", "c": "d", "e": "FG", "alpha": "", "beta": ""},
			[]string{"g"},
		},
		{
			Spec{"a": NoValueAccepted, "b": NoValueAccepted, "c": NoValueAccepted, "delta": ValueRequired},
			[]string{"-abc", "--", "ALPHA", "--beta"},
			OptionMap{"a": "", "b": "", "c": ""},
			[]string{"ALPHA", "--beta"},
		},
		{
			Spec{"a": NoValueAccepted, "b": NoValueAccepted, "c": ValueRequired, "alpha": ValueOptional},
			[]string{"-abc", "d", "--alpha", "BETA", "-efg", "--gamma", "g"},
			OptionMap{"a": "", "b": "", "c": "d", "alpha": ""},
			[]string{"BETA", "-efg", "--gamma", "g"},
		},
		{
			Spec{"1": NoValueAccepted, "a": ValueRequired},
			[]string{"-a", "=", "-1"},
			OptionMap{"1": "", "a": "="},
			[]string{},
		},
	}

	for _, tt := range testCases {
		t.Run(
			fmt.Sprintf("parses options from %q correctly", tt.input),
			func(t *testing.T) {
				t.Parallel()

				options, remainingArguments, err := ParseArgs(tt.input, tt.spec)
				if err != nil {
					t.Error(err)
				}

				if reflect.DeepEqual(options, tt.expectedOptions) == false {
					t.Errorf("expected %#v, got %#v", tt.expectedOptions, options)
				}

				if reflect.DeepEqual(remainingArguments, tt.expectedRemainingArguments) == false {
					t.Errorf("expected %#v, got %#v", tt.expectedRemainingArguments, remainingArguments)
				}
			},
		)
	}
}
