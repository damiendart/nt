package commands

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/damiendart/nt/internal/cli"
)

// New is a nt command that opens a new or existing note in Vim.
type New struct{}

// Run will execute the New command.
func (cmd *New) Run(app cli.Application, normalisedArgs []string) error {
	var argsEnd int
	var zettel bool

	for _, arg := range normalisedArgs {
		switch {
		case arg == "-z", arg == "--zettel":
			zettel = true

		default:
			normalisedArgs[argsEnd] = arg
			argsEnd++
		}
	}

	normalisedArgs = normalisedArgs[:argsEnd]

	if len(normalisedArgs) == 0 {
		return fmt.Errorf("new: missing filename")
	}

	file := normalisedArgs[0]

	if !strings.HasSuffix(file, ".md") {
		file = file + ".md"
	}

	if zettel {
		file = fmt.Sprintf(
			"%s/%s/%s/%s-%s",
			filepath.Dir(file),
			time.Now().Format("2006"),
			time.Now().Format("01"),
			time.Now().Format("200601021504"),
			filepath.Base(file),
		)
	}

	return cli.OpenFileInVim(app.Output, app.NotesDir, file)
}
