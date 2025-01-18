package main

import (
	"context"
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/ehrktia/lendbook/infra"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	log := JsonLogger(os.Stderr)
	if err := infra.Run(ctx, log); err != nil {
		log.LogAttrs(ctx, slog.LevelError, "stoping app",
			slog.String("error", err.Error()))
		os.Exit(1)
	}
}

func JsonLogger(w io.Writer) *slog.Logger {
	h := slog.NewJSONHandler(w, &slog.HandlerOptions{AddSource: true}).
		WithGroup("lendbook").WithAttrs([]slog.Attr{
		slog.String("version", "local"),
		slog.String("commit", "local"),
		slog.String("start-time", time.Now().UTC().String()),
	})
	l := slog.New(h)
	return l
}
