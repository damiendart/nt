// Copyright (C) Damien Dart, <damiendart@pobox.com>.
// This file is distributed under the MIT licence. For more information,
// please refer to the accompanying "LICENCE" file.

package slugify

import (
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

// Slugify returns a URL-friendly version of the string s.
func Slugify(s string) string {
	buf := make([]rune, 0, len(s))
	dash := false

	for _, r := range norm.NFKC.String(s) {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			buf = append(buf, unicode.ToLower(r))
			dash = true
		} else if dash && (unicode.IsMark(r) || unicode.IsPunct(r) || unicode.IsSpace(r) || unicode.IsSymbol(r)) {
			// Contractions and possessive forms like "don't" and
			// "Brown's" look better without the inner dash.
			if r == '\'' || r == '’' {
				continue
			}

			buf = append(buf, '-')
			dash = false
		}
	}

	return strings.Trim(string(buf), "-")
}
