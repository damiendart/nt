package editor

import (
	"io"
)

// Editor is implemented by anything that has a OpenFile method, which
// is used to open a file in a text editor.
type Editor interface {
	OpenFile(w io.Writer, file string) error
}
