package main

import (
	_ "embed"
	"fmt"

	"golang.org/x/tools/txtar"
)

//go:embed help.txtar
var src []byte

// Help represents a collection of help text files.
type Help struct {
	Texts map[string][]byte
}

// NewHelp returns a new instance of Help using text files from the
// accompanying txtar archive.
func NewHelp() *Help {
	a := txtar.Parse(src)
	h := &Help{Texts: make(map[string][]byte)}

	for _, f := range a.Files {
		h.Texts[f.Name] = f.Data
	}

	return h
}

// Get returns the help text associated with the provided key.
func (h *Help) Get(key string) ([]byte, error) {
	if t, found := h.Texts[key]; found {
		return t, nil
	}

	return nil, fmt.Errorf("unable to find help file %q", key)
}
