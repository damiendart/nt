package commands

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/damiendart/nt/internal/cli"
)

// Tags is a nt command that lists all tags used across all notes.
type Tags struct{}

// Run will execute the Tags command.
func (cmd *Tags) Run(app cli.Application, normalisedArgs []string) error {
	var showCount bool

	for _, arg := range normalisedArgs {
		if arg == "-c" {
			showCount = true
			break
		}

		if strings.HasPrefix(arg, "-") {
			return fmt.Errorf("tags: invalid option: %q", arg)
		}
	}

	tags := make(map[string]int)

	err := filepath.WalkDir(
		app.NotesDir,
		func(s string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			tagRegex, e := regexp.Compile("[ '(]#([a-zA-Z/:-]+)")
			if e != nil {
				return err
			}

			if filepath.Ext(s) == ".md" {
				f, err := os.Open(s)
				if err != nil {
					return err
				}
				defer f.Close()

				scanner := bufio.NewScanner(f)
				for scanner.Scan() {
					matches := tagRegex.FindAllStringSubmatch(scanner.Text(), -1)
					for _, m := range matches {
						tags[m[1]]++
					}
				}

				if err := scanner.Err(); err != nil {
					return err
				}
			}

			return nil
		},
	)
	if err != nil {
		return err
	}

	ts := make([]string, 0, len(tags))
	for k, v := range tags {
		if showCount {
			ts = append(ts, fmt.Sprintf("%s (%d)", k, v))
		} else {
			ts = append(ts, k)
		}
	}
	slices.Sort(ts)

	_, err = fmt.Fprintln(app.Output, strings.Join(ts, "\n"))

	return err
}
