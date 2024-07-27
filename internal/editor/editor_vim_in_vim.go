// Copyright (C) Damien Dart, <damiendart@pobox.com>.
// This file is distributed under the MIT licence. For more information,
// please refer to the accompanying "LICENCE" file.

package editor

import (
	"fmt"
	"io"
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

// OpenFile opens a file in current instance of Vim using Vim's terminal
// JSON API (for more information, search for "terminal-api" in Vim's
// help files).
func (e VimInVimEditor) OpenFile(path string, _ string) error {
	_, err := fmt.Fprintf(
		e.output,
		"\033]51;[%q, %q]\007",
		"drop",
		path,
	)

	return err
}
