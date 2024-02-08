package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/damiendart/nt/internal/cli"
)

// NewCommand is a nt command to open and create notes in a text editor.
type NewCommand struct{}

// Run will execute the NewCommand command.
func (cmd *NewCommand) Run(app Application, args []string) error {
	var zettel bool

	opts, remainingArgs, err := cli.ParseArgs(args, cli.Spec{"z": cli.NoValueAccepted})
	if err != nil {
		return err
	}

	for k := range opts {
		switch {
		case k == "z":
			zettel = true
		}
	}

	if len(remainingArgs) == 0 {
		return fmt.Errorf("new: missing filename")
	}

	file := remainingArgs[0]

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

	file = filepath.Join(app.NotesDir, file)
	file = filepath.Clean(file)

	if !strings.HasPrefix(file, app.NotesDir) {
		return fmt.Errorf("note would be created outside of %q", app.NotesDir)
	}

	err = os.MkdirAll(filepath.Dir(file), 0700)
	if err != nil {
		return err
	}

	err = os.Chdir(app.NotesDir)
	if err != nil {
		return err
	}

	return app.Editor.OpenFile(app.Output, file)
}
