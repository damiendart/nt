package commands

import (
	"fmt"
	"strings"

	"github.com/damiendart/nt/internal/cli"
)

type Inbox struct{}

func (cmd *Inbox) Run(app cli.Application, normalisedArgs []string) error {
	for _, arg := range normalisedArgs {
		if strings.HasPrefix(arg, "-") {
			return fmt.Errorf("invalid option: \"" + arg + "\"")
		}
	}

	return cli.OpenFileInVim(app.Output, app.NotesDir, "inbox.md")
}
