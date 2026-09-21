// Copyright (C) Damien Dart, <damiendart@pobox.com>.
// This file is distributed under the MIT licence. For more information,
// please refer to the accompanying "LICENCE" file.

package main

import (
	"errors"
	"net/http"
	"os"

	"github.com/yuin/goldmark-meta/v2"
	"github.com/yuin/goldmark/v2/extension"
	"github.com/yuin/goldmark/v2/parser"
	"github.com/yuin/goldmark/v2/renderer/html"

	"github.com/damiendart/nt/internal/cli"
	"github.com/damiendart/nt/internal/server"
)

// ServeCommand is a nt command to start a web server for viewing notes.
type ServeCommand struct{}

// Run will execute the InboxCommand command.
func (cmd *ServeCommand) Run(app Application, args []string) error {
	options, _, err := cli.ParseArgs(
		args,
		cli.Spec{
			"?":    cli.ValueOptional,
			"h":    cli.ValueOptional,
			"help": cli.ValueOptional,
		},
	)
	if err != nil {
		return err
	}

	for k := range options {
		switch {
		case k == "?", k == "h", k == "help":
			help, err := app.Help.Get("serve.txt")
			if err != nil {
				return err
			}

			_, err = app.Output.Write(help)
			if err != nil {
				return err
			}

			os.Exit(0)
		}
	}

	p := parser.New(
		parser.WithExtensions(
			extension.NewGFMParser(),
			extension.NewFootnoteParser(),
			extension.NewTypographerParser(),
			meta.Parser,
		),
	)
	r := html.New(
		html.WithExtensions(
			extension.NewGFMHTMLRenderer(),
			extension.NewFootnoteHTMLRenderer(),
		),
	)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /{name...}", server.NewNotesShowHandler(app.NotesRoot, p, r))

	srv := http.Server{
		Addr:    ":5555",
		Handler: mux,
	}

	err = srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
