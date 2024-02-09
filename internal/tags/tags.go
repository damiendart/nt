// Copyright (C) Damien Dart, <damiendart@pobox.com>.
// This file is distributed under the MIT licence. For more information,
// please refer to the accompanying "LICENCE" file.

package tags

import "regexp"

// If updating the following regular expression, the regular expressions
// in <https://github.com/damiendart/toolbox> to add tag syntax
// highlighting in Vim may also require updating.
var hashtagRegex = regexp.MustCompile("(?:^|[ \"'(])#([0-9/:-]*[a-zA-Z][a-zA-Z0-9/:-]*)")

// ExtractTags returns all tags in the provided string. Currently only
// hashtags are supported. Hashtags are phrases that start with a hash
// symbol and contain alphanumeric, "/", ":", and "-" characters.
// Hashtags must contain at least one letter.
func ExtractTags(s string) []string {
	var matches []string

	tagMatches := hashtagRegex.FindAllStringSubmatch(s, -1)
	for _, m := range tagMatches {
		matches = append(matches, m[1])
	}

	return matches
}
