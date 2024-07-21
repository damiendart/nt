// Copyright (C) Damien Dart, <damiendart@pobox.com>.
// This file is distributed under the MIT licence. For more information,
// please refer to the accompanying "LICENCE" file.

package editor

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
)

// VimEditor is a representation of the Vim text editor.
type VimEditor struct{}

// OpenFile opens a file in Vim. If in a Vim terminal, Vim's terminal
// JSON API will be used to open the file in the current instance of
// Vim (for more information, see "terminal-api" in Vim's help files).
func (editor VimEditor) OpenFile(path string, w io.Writer, root string) error {
	if _, ok := os.LookupEnv("VIM_TERMINAL"); ok {
		_, err := fmt.Fprintf(
			w,
			"\033]51;[%q, %q]\007",
			"drop",
			path,
		)

		return err
	}

	if root != "" {
		err := os.Chdir(root)
		if err != nil {
			return err
		}
	}

	cmd := exec.Command("vim", path)
	cmd.Env = os.Environ()
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	err := cmd.Run()
	if err != nil {
		var exitError *exec.ExitError

		if errors.As(err, &exitError) {
			return fmt.Errorf("vim exited with exit status %d", exitError.ExitCode())
		}

		return err
	}

	return nil
}
