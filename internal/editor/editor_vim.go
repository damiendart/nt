package editor

import (
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
func (editor VimEditor) OpenFile(w io.Writer, file string) error {
	if _, ok := os.LookupEnv("VIM_TERMINAL"); ok {
		_, err := fmt.Fprintf(
			w,
			"\033]51;[%q, %q]\007",
			"drop",
			file,
		)

		return err
	}

	cmd := exec.Command("vim", file)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	return cmd.Run()
}
