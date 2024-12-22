// Copyright (C) Damien Dart, <damiendart@pobox.com>.
// This file is distributed under the MIT licence. For more information,
// please refer to the accompanying "LICENCE" file.

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/damiendart/nt/internal/cli"
)

// JotCommand is a nt command that appends text with a timestamped
// heading to the top-level inbox note.
type JotCommand struct{}

// Run will execute the JotCommand command.
func (cmd JotCommand) Run(app Application, args []string) error {
	var text string

	options, _, err := cli.ParseArgs(
		args,
		cli.Spec{
			"?":    cli.ValueOptional,
			"h":    cli.ValueOptional,
			"help": cli.ValueOptional,
		},
	)
	if err != nil {
		return err
	}

	for k := range options {
		switch {
		case k == "?", k == "h", k == "help":
			help, err := app.Help.Get("jot.txt")
			if err != nil {
				return err
			}

			_, err = app.Output.Write(help)
			if err != nil {
				return err
			}

			os.Exit(0)
		}
	}

	if len(args) > 0 {
		text = strings.Join(args, " ")
	} else {
		fi, _ := os.Stdin.Stat()

		// Provide a simple prompt if no input has been piped in.
		if (fi.Mode() & os.ModeCharDevice) != 0 {
			_, err = fmt.Fprint(app.Output, "> ")
			if err != nil {
				return err
			}

			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()

			err = scanner.Err()
			if err != nil {
				return err
			}

			text = scanner.Text()
		} else {
			stdin, err := io.ReadAll(os.Stdin)
			if err != nil {
				return err
			}

			text = string(stdin)
		}
	}

	text = strings.TrimSuffix(text, "\n")

	if text == "" {
		app.Logger.Print("skipping empty input")

		return nil
	}

	f, err := os.OpenFile(
		filepath.Join(app.NotesDir, "inbox.md"),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0600,
	)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(f, "\n### %s\n\n%s\n", time.Now().Format(time.RFC1123), text)
	if err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	return nil
}
