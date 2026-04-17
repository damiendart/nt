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
	"strings"

	"github.com/damiendart/nt/internal/cli"
)

// BacklinksCommand is a command to list backlinks to a particular note.
type BacklinksCommand struct{}

// Run will execute the BacklinksCommand command.
func (cmd *BacklinksCommand) Run(app Application, args []string) error {
	opts, remainingArgs, err := cli.ParseArgs(
		args,
		cli.Spec{
			"?":    cli.ValueOptional,
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
			help, err := app.Help.Get("backlinks.txt")
			if err != nil {
				return err
			}

			_, err = app.Output.Write(help)
			if err != nil {
				return err
			}

			os.Exit(0)
		}
	}

	if len(remainingArgs) == 0 {
		return fmt.Errorf("missing filename")
	}

	base := strings.Join(remainingArgs[0:], " ")
	backlinksFound := false

	if filepath.Ext(base) == ".md" {
		base = base[0 : len(base)-len(filepath.Ext(base))]
	}

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

				candidates := []string{
					"[[" + base + "]]",
					"[[" + base + "|",
					"[[" + base + ".md]]",
					"[[" + base + ".md|",
				}

				if r, err := filepath.Rel(filepath.Dir(s), base); err == nil {
					candidates = append(candidates, "[["+r+"]]")
					candidates = append(candidates, "[["+r+"|")
					candidates = append(candidates, "[["+r+".md]]")
					candidates = append(candidates, "[["+r+".md|")
				}

				line := 0
				scanner := bufio.NewScanner(f)
				for scanner.Scan() {
					line++

					for _, candidate := range candidates {
						if i := strings.Index(scanner.Text(), candidate); i != -1 {
							backlinksFound = true
							_, err := fmt.Fprintf(app.Output, "%s:%d:%d:%s\n", s, line, i+3, scanner.Text())
							if err != nil {
								return err
							}
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

	if !backlinksFound {
		return NoResultsError("no backlinks found")
	}

	return nil
}
