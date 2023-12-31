package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// An Application is used to store any application-wide dependencies.
type Application struct {
	Commands map[string]Command
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

func normaliseArgs(args []string) []string {
	var normalisedArgs []string

	for _, arg := range args {
		if strings.HasPrefix(arg, "--") {
			normalisedArgs = append(
				normalisedArgs,
				strings.Split(arg, "=")...,
			)
		} else if strings.HasPrefix(arg, "-") {
			for i, r := range arg[1:] {
				if r == '=' {
					normalisedArgs = append(normalisedArgs, arg[2+i:])

					break
				}

				normalisedArgs = append(normalisedArgs, "-"+string(r))
			}
		} else {
			normalisedArgs = append(normalisedArgs, arg)
		}
	}

	return normalisedArgs
}

func main() {
	var argsEnd int
	var currentOption string
	var notesDir string

	cmdMap := map[string]Command{
		"inbox": &InboxCommand{},
		"jot":   &JotCommand{},
		"new":   &NewCommand{},
		"tags":  &TagsCommand{},
	}
	logger := log.New(os.Stderr, os.Args[0]+": ", 0)
	normalisedArgs := normaliseArgs(os.Args[1:])

	// Perform an initial pass through the command line argument list to
	// handle any global options.
	for _, arg := range normalisedArgs {
		if currentOption != "" && strings.HasPrefix(arg, "-") {
			logger.Fatalf("Missing value for %q option", currentOption)
		}

		switch {
		case arg == "-?", arg == "-h", arg == "--help":
			fmt.Println(help(cmdMap))
			os.Exit(0)

		case arg == "--version":
			fmt.Println("VERSION GOES HERE")
			os.Exit(0)

		case arg == "--notes-dir":
			currentOption = arg

		default:
			switch currentOption {
			case "--notes-dir":
				notesDir = arg

			default:
				normalisedArgs[argsEnd] = arg
				argsEnd++
			}

			currentOption = ""
		}
	}

	normalisedArgs = normalisedArgs[:argsEnd]

	if len(normalisedArgs) == 0 {
		logger.Fatalf("missing command")
	}

	if currentOption != "" {
		logger.Fatalf("missing value for %q option", currentOption)
	}

	if strings.HasPrefix(normalisedArgs[0], "-") {
		logger.Fatalf("invalid option: \"" + normalisedArgs[0] + "\"")
	}

	if notesDir == "" {
		if env := os.Getenv("NOTES_ROOT"); env != "" {
			notesDir = filepath.Clean(env)
		} else {
			logger.Fatalf("notes directory not set")
		}
	}

	application := &Application{
		Commands: cmdMap,
		Logger:   logger,
		NotesDir: notesDir,
		Output:   os.Stdout,
	}

	command, ok := application.Commands[normalisedArgs[0]]
	if !ok {
		application.Logger.Fatalf("invalid command: %q", normalisedArgs[0])
	}

	err := command.Run(*application, normalisedArgs[1:])
	if err != nil {
		application.Logger.Fatal(err)
	}
}
