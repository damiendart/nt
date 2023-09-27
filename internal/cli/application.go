package cli

import (
	"io"
	"log"
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
