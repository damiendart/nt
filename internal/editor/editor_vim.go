// Copyright (C) Damien Dart, <damiendart@pobox.com>.
// This file is distributed under the MIT licence. For more information,
// please refer to the accompanying "LICENCE" file.

package editor

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
)

// VimEditor is a representation of the Vim text editor.
type VimEditor struct{}

// Open implements the editor.Opener interface. If cwd is provided and
// not an empty string, the working directory will be updated.
func (editor VimEditor) Open(path string, cwd string) error {
	if !filepath.IsAbs(path) {
		return ErrPathNotAbsolute
	}

	// Windows does not support the "execve(2)" system call.
	if runtime.GOOS == "windows" {
		if cwd != "" {
			err := os.Chdir(cwd)
			if err != nil {
				return err
			}
		}

		cmd := exec.Command("vim", quotePath(path))
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

	shell := os.Getenv("SHELL")
	if len(shell) == 0 {
		shell = "sh"
	}

	shellPath, err := exec.LookPath(shell)
	if err != nil {
		return err
	}

	var cmd string

	// Standard input is reopened as "/dev/tty" to prevent Vim from
	// bombing out if input has been piped in. Currently, any piped
	// input is not passed through to Vim.
	if cwd == "" {
		cmd = fmt.Sprintf("vim %s </dev/tty", quotePath(path))
	} else {
		cmd = fmt.Sprintf(`vim --cmd 'cd %s' %s </dev/tty`, cwd, quotePath(path))
	}

	return syscall.Exec(shellPath, []string{shellPath, "-c", cmd}, os.Environ())
}

func quotePath(path string) string {
	r := strings.NewReplacer("'", "'\\''")

	return "'" + r.Replace(path) + "'"
}
