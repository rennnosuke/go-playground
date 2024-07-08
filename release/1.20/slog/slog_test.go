package slog

import (
	"os"
	"testing"

	"golang.org/x/exp/slog"
)

func TestInfo(t *testing.T) {
	slog.Info("hello world!", "name", "blue", "No", 5)
}

func TestJSONHandler(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout))
	logger.Info("json log", "key1", "value1")
}

func TestJSONHandlerWithType(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout))
	logger.Info("json log with type",
		slog.String("str", "hogehoge"),
		slog.Int("int", 42),
		slog.Bool("bool", true),
		slog.Any("nested1", struct {
			Name string
		}{
			Name: "name",
		}),
		slog.Group("nested2",
			slog.String("inner_str", "hogehoge"),
		),
	)
}

func TestAddSource(t *testing.T) {
	logger := slog.New(slog.HandlerOptions{AddSource: true}.NewJSONHandler(os.Stdout))
	logger.Info("json log", "key1", "value1")
}

func TestLoggerWith(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout))
	logger2 := logger.With(
		slog.String("method", "GET"),
	)
	logger2.Info("json log", "key1", "value1")
}
