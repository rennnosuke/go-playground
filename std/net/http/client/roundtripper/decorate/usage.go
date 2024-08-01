package main

import (
	"log/slog"
	"net/http"
)

var cli = http.Client{
	Transport: &DecorateRoundTripper{
		base: http.DefaultTransport,
		before: func(r *http.Request) error {
			slog.InfoContext(r.Context(), "before request")
			return nil
		},
		after: func(r *http.Response) error {
			slog.InfoContext(r.Request.Context(), "after response")
			return nil
		},
	},
	Timeout: 30,
}
