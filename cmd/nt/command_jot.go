package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// JotCommand is a nt command that appends a timestamped Markdown list
// item to the top-level inbox note.
type JotCommand struct{}

// Run will execute the JotCommand command.
func (cmd JotCommand) Run(app Application, normalisedArgs []string) error {
	var text string

	for _, arg := range normalisedArgs {
		if strings.HasPrefix(arg, "-") {
			return fmt.Errorf("invalid option: \"" + arg + "\"")
		}
	}

	if len(normalisedArgs) > 0 {
		text = strings.Join(normalisedArgs, " ")

		defer func() {
			fmt.Fprintln(app.Output, text)
		}()
	} else {
		_, err := fmt.Fprint(app.Output, "> ")
		if err != nil {
			return err
		}

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()

		err = scanner.Err()
		if err != nil {
			return err
		}

		text = scanner.Text()
	}

	f, err := os.OpenFile(
		filepath.Join(app.NotesDir, "inbox.md"),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0600,
	)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = fmt.Fprintf(f, "- %s: %s\n", time.Now().Format(time.RFC1123), text)
	if err != nil {
		return err
	}

	return nil
}
