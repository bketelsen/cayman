package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"cayman/internal/modules"
)

func main() {
	var (
		addr = flag.String("addr", "0.0.0.0", "listen address")
		port = flag.String("port", "8080", "listen port")
	)
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})
	logger := slog.New(handler)
	slog.SetDefault(logger)

	engine := modules.NewEngine(logger, *addr, *port)
	if err := engine.Start(ctx); err != nil {
		logger.Error("failed to start engine", "error", err)
	}
}
