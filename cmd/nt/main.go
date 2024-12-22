// Copyright (C) Damien Dart, <damiendart@pobox.com>.
// This file is distributed under the MIT licence. For more information,
// please refer to the accompanying "LICENCE" file.

package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/damiendart/nt/internal/cli"
	"github.com/damiendart/nt/internal/editor"
)

// An Application is used to store any application-wide dependencies.
type Application struct {
	Commands map[string]Command
	Editor   editor.Opener
	Logger   *log.Logger
	NotesDir string
	Output   io.Writer
}

// Command is implemented by anything that has a Run method. The
// implementation can then be used as nt command.
type Command interface {
	Run(app Application, normalisedArgs []string) error
}

func help(cmdMap map[string]Command) string {
	return `HELP TEXT GOES HERE`
}

func main() {
	var notesDir string

	cmdMap := map[string]Command{
		"inbox": &InboxCommand{},
		"jot":   &JotCommand{},
		"new":   &NewCommand{},
		"tags":  &TagsCommand{},
	}
	logger := log.New(os.Stderr, os.Args[0]+": ", 0)

	globalOptions, remainingArgs, err := cli.ParseArgs(
		os.Args[1:],
		cli.Spec{
			"?":         cli.ValueOptional,
			"h":         cli.ValueOptional,
			"help":      cli.ValueOptional,
			"notes-dir": cli.ValueRequired,
			"version":   cli.NoValueAccepted,
		},
	)
	if err != nil {
		logger.Fatalf(err.Error())
	}

	for k, v := range globalOptions {
		switch {
		case k == "?", k == "h", k == "help":
			fmt.Println(help(cmdMap))
			os.Exit(0)

		case k == "version":
			fmt.Println("VERSION GOES HERE")
			os.Exit(0)

		case k == "notes-dir":
			notesDir = v
		}
	}

	if len(remainingArgs) == 0 {
		logger.Fatalf("missing command")
	}

	if notesDir == "" {
		if env := os.Getenv("NOTES_ROOT"); env != "" {
			notesDir = filepath.Clean(env)
		} else {
			logger.Fatalf("notes directory not set")
		}
	}

	var e editor.Opener

	_, ok := os.LookupEnv("VIM_TERMINAL")
	if ok {
		e = editor.NewVimInVimEditor(os.Stdout)
	} else {
		e = &editor.VimEditor{}
	}

	application := &Application{
		Commands: cmdMap,
		Editor:   e,
		Logger:   logger,
		NotesDir: notesDir,
		Output:   os.Stdout,
	}

	command, ok := application.Commands[remainingArgs[0]]
	if !ok {
		application.Logger.Fatalf("error: invalid command: %q", remainingArgs[0])
	}

	application.Logger.SetPrefix(
		fmt.Sprintf(
			"%s%s: ",
			application.Logger.Prefix(),
			remainingArgs[0],
		),
	)

	err = command.Run(*application, remainingArgs[1:])
	if err != nil {
		application.Logger.Fatalf("error: %s", err)
	}
}
