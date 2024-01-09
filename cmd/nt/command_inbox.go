package main

import (
	"fmt"
	"os"
	"strings"
)

// InboxCommand is a nt command to open the top-level inbox note in a
// text editor.
type InboxCommand struct{}

// Run will execute the InboxCommand command.
func (cmd *InboxCommand) Run(app Application, normalisedArgs []string) error {
	for _, arg := range normalisedArgs {
		if strings.HasPrefix(arg, "-") {
			return fmt.Errorf("invalid option: \"" + arg + "\"")
		}
	}

	err := os.Chdir(app.NotesDir)
	if err != nil {
		return err
	}

	return app.Editor.OpenFile(app.Output, "inbox.md")
}
