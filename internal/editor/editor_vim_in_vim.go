// Copyright (C) Damien Dart, <damiendart@pobox.com>.
// This file is distributed under the MIT licence. For more information,
// please refer to the accompanying "LICENCE" file.

package editor

import (
	"fmt"
	"io"
	"path/filepath"
)

// VimInVimEditor is a representation of the Vim text editor being
// opened in a Vim terminal.
type VimInVimEditor struct {
	output io.Writer
}

// NewVimInVimEditor returns a new instance of VimInVimEditor.
func NewVimInVimEditor(output io.Writer) *VimInVimEditor {
	return &VimInVimEditor{output: output}
}

// Open implements the editor.Opener interface. Open uses Vim's terminal
// JSON API to open a file specified at path in the current instance of
// Vim (for more information, search for "terminal-api" in Vim's help).
func (e VimInVimEditor) Open(path string, _ string) error {
	if !filepath.IsAbs(path) {
		return ErrPathNotAbsolute
	}

	_, err := fmt.Fprintf(
		e.output,
		"\033]51;[%q, %q]\007",
		"drop",
		path,
	)

	return err
}
