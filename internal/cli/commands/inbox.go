package commands

import (
	"fmt"
	"strings"

	"github.com/damiendart/nt/internal/cli"
)

// Inbox is a nt command that opens a top-level inbox note in Vim.
type Inbox struct{}

// Run will execute the Inbox command.
func (cmd *Inbox) Run(app cli.Application, normalisedArgs []string) error {
	for _, arg := range normalisedArgs {
		if strings.HasPrefix(arg, "-") {
			return fmt.Errorf("invalid option: \"" + arg + "\"")
		}
	}

	return cli.OpenFileInVim(app.Output, app.NotesDir, "inbox.md")
}
