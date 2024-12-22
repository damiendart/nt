// Copyright (C) Damien Dart, <damiendart@pobox.com>.
// This file is distributed under the MIT licence. For more information,
// please refer to the accompanying "LICENCE" file.

package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/damiendart/nt/internal/cli"
	"github.com/damiendart/nt/internal/tags"
)

// TagsCommand is a nt command to list all tags used across all notes.
type TagsCommand struct{}

// Run will execute the TagsCommand command.
func (cmd *TagsCommand) Run(app Application, args []string) error {
	var showCount bool

	opts, _, err := cli.ParseArgs(
		args,
		cli.Spec{
			"?":    cli.ValueOptional,
			"c":    cli.NoValueAccepted,
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
		}
	}

	tagsCount := make(map[string]int)

	err = filepath.WalkDir(
		app.NotesDir,
		func(s string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if filepath.Ext(s) == ".md" {
				f, err := os.Open(s)
				if err != nil {
					return err
				}
				defer f.Close()

				scanner := bufio.NewScanner(f)
				for scanner.Scan() {
					matches := tags.ExtractHashtags(scanner.Text())
					for _, t := range matches {
						tagsCount[t]++
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

	ts := make([]string, 0, len(tagsCount))
	for k, v := range tagsCount {
		if showCount {
			ts = append(ts, fmt.Sprintf("%s (%d)", k, v))
		} else {
			ts = append(ts, k)
		}
	}
	slices.Sort(ts)

	_, err = fmt.Fprintln(app.Output, strings.Join(ts, "\n"))

	return err
}
