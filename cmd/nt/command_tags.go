// Copyright (C) Damien Dart, <damiendart@pobox.com>.
// This file is distributed under the MIT licence. For more information,
// please refer to the accompanying "LICENCE" file.

package main

import (
	"bufio"
	"cmp"
	"fmt"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/damiendart/nt/internal/cli"
	"github.com/damiendart/nt/internal/fuzzy"
	"github.com/damiendart/nt/internal/tags"
)

// TagsCommand is a nt command to list all tags used across all notes.
type TagsCommand struct{}

type match struct {
	count int
	score int
	tag   string
}

// Run will execute the TagsCommand command.
func (cmd *TagsCommand) Run(app Application, args []string) error {
	var showCount bool
	var useFuzzyMatching bool

	opts, remainingArgs, err := cli.ParseArgs(
		args,
		cli.Spec{
			"?":    cli.ValueOptional,
			"c":    cli.NoValueAccepted,
			"f":    cli.NoValueAccepted,
			"h":    cli.ValueOptional,
			"help": cli.ValueOptional,
		},
	)
	if err != nil {
		return err
	}

	for k := range opts {
		switch {
		case k == "?", k == "h", k == "help":
			help, err := app.Help.Get("jot.txt")
			if err != nil {
				return err
			}

			_, err = app.Output.Write(help)
			if err != nil {
				return err
			}

			os.Exit(0)

		case k == "c":
			showCount = true

		case k == "f":
			useFuzzyMatching = true
		}
	}

	matches := make(map[string]match)

	err = fs.WalkDir(
		app.NotesRoot.FS(),
		".",
		func(s string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if filepath.Ext(s) == ".md" {
				f, err := app.NotesRoot.Open(s)
				if err != nil {
					return err
				}
				defer f.Close()

				scanner := bufio.NewScanner(f)
				for scanner.Scan() {
					tags := tags.ExtractHashtags(scanner.Text())
					for _, tag := range tags {
						if m, ok := matches[tag]; ok {
							m.count++
							matches[tag] = m

							continue
						}

						if len(remainingArgs) > 0 {
							if useFuzzyMatching {
								if s := fuzzy.MatchScore(tag, remainingArgs[0]); s != -1 {
									matches[tag] = match{1, s, tag}
								}
							} else if strings.HasPrefix(tag, remainingArgs[0]) {
								matches[tag] = match{1, 0, tag}
							}
						} else {
							matches[tag] = match{1, 0, tag}
						}
					}
				}

				if err := scanner.Err(); err != nil {
					return err
				}
			}

			return nil
		},
	)
	if err != nil {
		return err
	}

	ms := slices.Collect(maps.Values(matches))

	slices.SortFunc(
		ms,
		func(a, b match) int {
			if n := cmp.Compare(a.score, b.score); n != 0 {
				return n
			}
			return strings.Compare(a.tag, b.tag)
		},
	)

	if len(ms) > 0 {
		for _, m := range ms {
			if showCount {
				_, err := fmt.Fprintf(app.Output, "%s (%d)\n", m.tag, m.count)
				if err != nil {
					return err
				}
			} else {
				_, err := fmt.Fprintf(app.Output, "%s\n", m.tag)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
