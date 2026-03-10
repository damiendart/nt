// Copyright (C) Damien Dart, <damiendart@pobox.com>.
// This file is distributed under the MIT licence. For more information,
// please refer to the accompanying "LICENCE" file.

package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/damiendart/nt/internal/cli"
	"github.com/damiendart/nt/internal/slugify"
)

// NewCommand is a nt command to open and create notes in a text editor.
type NewCommand struct{}

// Run will execute the NewCommand command.
func (cmd *NewCommand) Run(app Application, args []string) error {
	var createDirs bool
	var zettel bool

	opts, remainingArgs, err := cli.ParseArgs(
		args,
		cli.Spec{
			"?":           cli.ValueOptional,
			"create-dirs": cli.NoValueAccepted,
			"d":           cli.NoValueAccepted,
			"h":           cli.ValueOptional,
			"help":        cli.ValueOptional,
			"z":           cli.NoValueAccepted,
		},
	)
	if err != nil {
		return err
	}

	for k := range opts {
		switch {
		case k == "?", k == "h", k == "help":
			help, err := app.Help.Get("new.txt")
			if err != nil {
				return err
			}

			_, err = app.Output.Write(help)
			if err != nil {
				return err
			}

			os.Exit(0)
		case k == "create-dirs", k == "d":
			createDirs = true
		case k == "z":
			zettel = true
		}
	}

	if len(remainingArgs) == 0 {
		return fmt.Errorf("missing filename")
	}

	file := strings.Join(remainingArgs[0:], " ")
	file = file[0 : len(file)-len(filepath.Ext(file))]

	if zettel {
		file = fmt.Sprintf(
			"%s/%s/%s/%s-%s.md",
			filepath.Dir(file),
			time.Now().Format("2006"),
			time.Now().Format("01"),
			time.Now().Format("200601021504"),
			slugify.Slugify(filepath.Base(file)),
		)
	} else {
		file = fmt.Sprintf(
			"%s/%s.md",
			filepath.Dir(file),
			slugify.Slugify(filepath.Base(file)),
		)
	}

	file = filepath.Clean(file)

	if createDirs {
		err := app.NotesRoot.MkdirAll(filepath.Dir(file), 0700)
		if err != nil {
			if pathError, ok := errors.AsType[*os.PathError](err); ok {
				// See <https://github.com/golang/go/issues/74640>.
				if strings.Contains(pathError.Error(), "path escapes from parent") {
					return fmt.Errorf("note would be created outside of %q", app.NotesRoot.Name())
				}

				return err
			}

			return err
		}
	} else {
		if !strings.HasPrefix(filepath.Join(app.NotesRoot.Name(), file), app.NotesRoot.Name()) {
			return fmt.Errorf("note would be created outside of %q", app.NotesRoot.Name())
		}
	}

	return app.Editor.Open(
		filepath.Join(app.NotesRoot.Name(), file),
		app.NotesRoot.Name(),
	)
}
