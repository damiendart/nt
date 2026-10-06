// Copyright (C) Damien Dart, <damiendart@pobox.com>.
// This file is distributed under the MIT licence. For more information,
// please refer to the accompanying "LICENCE" file.

package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

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
	var port uint = 5555

	options, _, err := cli.ParseArgs(
		args,
		cli.Spec{
			"?":    cli.ValueOptional,
			"h":    cli.ValueOptional,
			"p":    cli.ValueRequired,
			"help": cli.ValueOptional,
			"port": cli.ValueRequired,
		},
	)
	if err != nil {
		return err
	}

	for k, v := range options {
		switch {
		case k == "p", k == "port":
			p, err := strconv.ParseUint(v, 10, 32)
			if err != nil || p > 65535 {
				return fmt.Errorf("invalid port number %v", v)
			}

			port = uint(p)
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

	logger := slog.New(slog.NewJSONHandler(app.Output, &slog.HandlerOptions{Level: slog.LevelDebug}))
	mux := server.NewRouter()

	mux.Use(server.DefaultHeadersMiddleware, server.NewLogRequestMiddleware(logger))
	mux.HandleFunc("GET /{name...}", server.NewNotesShowHandler(app.NotesRoot, p, r))

	srv := http.Server{
		Addr:     fmt.Sprintf(":%d", port),
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelWarn),
		Handler:  mux,
	}

	shutdownErr := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		s := <-quit

		logger.LogAttrs(
			context.Background(),
			slog.LevelInfo,
			"stopping server",
			slog.String("addr", srv.Addr),
			slog.String("signal", s.String()),
		)

		shutdownErr <- srv.Shutdown(context.TODO())
	}()

	logger.LogAttrs(
		context.Background(),
		slog.LevelInfo,
		"starting server",
		slog.String("addr", srv.Addr),
	)

	err = srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownErr
	if err != nil {
		return err
	}

	return nil
}
