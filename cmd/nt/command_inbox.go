// Copyright (C) Damien Dart, <damiendart@pobox.com>.
// This file is distributed under the MIT licence. For more information,
// please refer to the accompanying "LICENCE" file.

package main

import (
	"os"
	"path/filepath"

	"github.com/damiendart/nt/internal/cli"
)

// InboxCommand is a nt command to open the top-level inbox note in a
// text editor.
type InboxCommand struct{}

// Run will execute the InboxCommand command.
func (cmd *InboxCommand) Run(app Application, args []string) error {
	options, _, err := cli.ParseArgs(
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

	for k := range options {
		switch {
		case k == "?", k == "h", k == "help":
			help, err := app.Help.Get("inbox.txt")
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

	return app.Editor.Open(
		filepath.Join(app.NotesDir, "inbox.md"),
		app.NotesDir,
	)
}
