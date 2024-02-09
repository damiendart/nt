// Copyright (C) Damien Dart, <damiendart@pobox.com>.
// This file is distributed under the MIT licence. For more information,
// please refer to the accompanying "LICENCE" file.

package main

import (
	"os"

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

	err = os.Chdir(app.NotesDir)
	if err != nil {
		return err
	}

	return app.Editor.OpenFile(app.Output, "inbox.md")
}
