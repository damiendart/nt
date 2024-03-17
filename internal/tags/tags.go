// Copyright (C) Damien Dart, <damiendart@pobox.com>.
// This file is distributed under the MIT licence. For more information,
// please refer to the accompanying "LICENCE" file.

package tags

import (
	"strings"
	"unicode"
)

// ExtractHashtags returns all hashtags in the provided string. Hashtags
// are words that start with a hash symbol and contain alphanumeric,
// "/", ":", and "-" characters. Hashtags must contain at least one
// letter, and any valid trailing punctuation characters are omitted.
// Hashtags can be surrounded with quotation marks and parentheses.
func ExtractHashtags(s string) []string {
	var matches []string

	for _, t := range strings.Fields(s) {
		t = strings.Trim(t, "\"'()")

		if len(t) == 1 {
			continue
		}

		if !strings.HasPrefix(t, "#") {
			continue
		}

		if isValidHashtagName(t[1:]) {
			matches = append(matches, strings.Trim(t[1:], "/:-"))
		}
	}

	return matches
}

func isValidHashtagName(s string) bool {
	hasLetter := false

	for _, r := range s {
		if unicode.IsLetter(r) {
			hasLetter = true
			continue
		}

		if unicode.IsNumber(r) {
			continue
		}

		if r == '/' || r == ':' || r == '-' {
			continue
		}

		return false
	}

	return hasLetter
}
