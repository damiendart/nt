// Copyright (C) Damien Dart, <damiendart@pobox.com>.
// This file is distributed under the MIT licence. For more information,
// please refer to the accompanying "LICENCE" file.

package server

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/yuin/goldmark/v2/parser"
	"github.com/yuin/goldmark/v2/renderer/html"
)

// NewNotesShowHandler returns a handler for showing rendered notes.
func NewNotesShowHandler(root *os.Root, p parser.Parser, h html.Renderer) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimSuffix(r.PathValue("name"), ".md") + ".md"

		md, err := root.Open(path)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		defer md.Close()

		content, err := io.ReadAll(md)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		ast := p.Parse(content)

		var buf bytes.Buffer
		if err := h.Render(&buf, content, ast); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(200)
		w.Write(buf.Bytes())
	}
}
