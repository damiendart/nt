package cli

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// OpenFileInVim opens a file in Vim.
func OpenFileInVim(w io.Writer, cwd string, filename string) error {
	file := filepath.Join(cwd, filename)

	if !strings.HasPrefix(file, cwd) {
		return fmt.Errorf("file would be created outside of %q", cwd)
	}

	err := os.MkdirAll(filepath.Dir(file), 0700)
	if err != nil {
		return err
	}

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
	cmd.Dir = cwd
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	return cmd.Run()
}
