// Copyright (C) Damien Dart, <damiendart@pobox.com>.
// This file is distributed under the MIT licence. For more information,
// please refer to the accompanying "LICENCE" file.

// Note: If anything here gets updated, the hashtag-related
// bits-and-pieces in my Vim configuration files (available at
// <https://github.com/damiendart/toolbox>) may also need updating.

package tags

import (
	"strings"
	"unicode"
)

// ExtractHashtags returns all hashtags (without starting hash symbols)
// in the provided string.
//
// Hashtags are words that start with a hash symbol, followed by
// alphanumeric, "/", ":", "-", and "_" characters. Hashtags must
// contain at least one letter and any trailing colons are ignored. They
// can be surrounded with quotation marks and parentheses.
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
			matches = append(matches, strings.TrimRight(t[1:], ":"))
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

		if r == '/' || r == ':' || r == '-' || r == '_' {
			continue
		}

		return false
	}

	return hasLetter
}
