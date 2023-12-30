package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/damiendart/nt/internal/editor"
)

// NewCommand is a nt command that opens a new or existing note in Vim.
type NewCommand struct{}

// Run will execute the NewCommand command.
func (cmd *NewCommand) Run(app Application, normalisedArgs []string) error {
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

	if !strings.HasPrefix(file, app.NotesDir) {
		return fmt.Errorf("note would be created outside of %q", app.NotesDir)
	}

	err := os.MkdirAll(filepath.Dir(file), 0700)
	if err != nil {
		return err
	}

	return editor.OpenFileInVim(app.Output, app.NotesDir, file)
}
