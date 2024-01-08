package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/damiendart/nt/internal/editor"
)

// InboxCommand is a nt command to open the top-level inbox note in Vim.
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

	return editor.OpenFileInVim(app.Output, "inbox.md")
}
