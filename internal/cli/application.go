package cli

import (
	"io"
	"log"
)

type Application struct {
	Commands map[string]Command
	Logger   *log.Logger
	NotesDir string
	Output   io.Writer
}

type Command interface {
	Run(app Application, normalisedArgs []string) error
}
