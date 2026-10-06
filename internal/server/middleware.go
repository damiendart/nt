// Copyright (C) Damien Dart, <damiendart@pobox.com>.
// This file is distributed under the MIT licence. For more information,
// please refer to the accompanying "LICENCE" file.

package server

import (
	"log/slog"
	"net/http"
	"time"
)

type metricsResponseWriter struct {
	http.ResponseWriter
	StatusCode    int
	BytesCount    int
	headerWritten bool
}

func (w *metricsResponseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)

	if !w.headerWritten {
		w.StatusCode = statusCode
		w.headerWritten = true
	}
}

func (w *metricsResponseWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)

	w.BytesCount += n
	w.headerWritten = true

	return n, err
}

func (w *metricsResponseWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}

// DefaultHeadersMiddleware is a HTTP middleware function that adds a
// few common HTTP headers that apply to all requests.
func DefaultHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("X-Content-Type-Options", "nosniff")
			w.Header().Add("X-Frame-Options", "deny")

			next.ServeHTTP(w, r)
		},
	)
}

// NewLogRequestMiddleware returns a request/response logging HTTP
// middleware function.
func NewLogRequestMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				mw := metricsResponseWriter{
					ResponseWriter: w,
					StatusCode:     http.StatusOK,
				}

				start := time.Now()
				next.ServeHTTP(&mw, r)
				duration := time.Since(start)

				level := slog.LevelInfo

				if mw.StatusCode == http.StatusInternalServerError {
					level = slog.LevelError
				}

				logger.LogAttrs(
					r.Context(),
					level,
					"access",
					slog.GroupAttrs(
						"request",
						slog.String("method", r.Method),
						slog.String("path", r.URL.String()),
						slog.String("proto", r.Proto),
					),
					slog.GroupAttrs(
						"response",
						slog.Int("status_code", mw.StatusCode),
						slog.Duration("duration_ns", duration),
						slog.Int("body_bytes", mw.BytesCount),
					),
				)
			},
		)
	}
}
