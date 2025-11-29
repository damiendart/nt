package fuzzy

import (
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

// IsFuzzyMatch reports whether the characters in the given needle
// string occur in the same order as the given haystack string. Ahead of
// the comparison, both needle and haystack strings are converted to
// lowercase and character accents are removed.
func IsFuzzyMatch(needle string, haystack string) bool {
	haystack = normalise(haystack)
	needle = normalise(needle)

	p := 0

	for _, r := range needle {
		i := strings.IndexRune(haystack[p:], r)

		if i == -1 {
			return false
		}

		p += i + 1
	}

	return true
}

func normalise(s string) string {
	return strings.Map(
		func(r rune) rune {
			if unicode.IsMark(r) || unicode.IsSpace(r) {
				return -1
			}

			return unicode.ToLower(r)
		},
		norm.NFD.String(s),
	)
}
