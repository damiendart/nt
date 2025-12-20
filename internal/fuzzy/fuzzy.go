package fuzzy

import (
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

// MatchScore returns the length of the shortest prefix of s that
// contains all characters of needle in order, but not necessarily
// sequentially. If no such prefix exists, MatchScore returns -1.
// When comparing match scores, lower positive values are better.
//
// Ahead of the comparison, both needle and s are converted to lowercase
// and any diacritical marks are removed.
func MatchScore(s string, needle string) int {
	var p int

	needle = normalise(needle)
	s = normalise(s)

	if len(needle) == 0 || len(needle) > len(s) {
		return -1
	}

	for _, r := range needle {
		i := strings.IndexRune(s[p:], r)

		if i == -1 {
			return -1
		}

		p += i + 1
	}

	return p
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
