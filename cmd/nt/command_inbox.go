// Copyright (C) Damien Dart, <damiendart@pobox.com>.
// This file is distributed under the MIT licence. For more information,
// please refer to the accompanying "LICENCE" file.

package main

import (
	"path/filepath"

	"github.com/damiendart/nt/internal/cli"
)

// InboxCommand is a nt command to open the top-level inbox note in a
// text editor.
type InboxCommand struct{}

// Run will execute the InboxCommand command.
func (cmd *InboxCommand) Run(app Application, args []string) error {
	_, _, err := cli.ParseArgs(args, cli.Spec{})
	if err != nil {
		return err
	}

	return app.Editor.OpenFile(
		filepath.Join(app.NotesDir, "inbox.md"),
		app.Output,
		app.NotesDir,
	)
}
