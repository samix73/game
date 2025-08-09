package main

import (
	"context"
	"flag"
	"log/slog"
	"os"

	"github.com/samix73/game/game"
	"github.com/samix73/game/worlds"
	"golang.org/x/image/math/f64"
)

var (
	fullscreen = flag.Bool("fullscreen", false, "enable fullscreen mode")
	tracing    = flag.Bool("tracing", false, "enable tracing")
	logLevel   = flag.String("log-level", "info", "set the log level (debug, info, warn, error, fatal)")
)

func setupLogger(level string) {
	var logLevel slog.Level
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo // Default to info if an invalid level is provided
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))
	slog.SetDefault(logger)
}

func main() {
	flag.Parse()

	setupLogger(*logLevel)

	ctx, _ := context.WithCancel(context.Background())

	g := game.NewGame(ctx, &game.Config{
		Title:        "Game",
		ScreenWidth:  1280,
		ScreenHeight: 960,
		Gravity:      f64.Vec2{0, 9.81},
		Fullscreen:   *fullscreen,
		Tracing:      *tracing,
	})

	mainWorld, err := worlds.NewMainWorld(ctx, g)
	if err != nil {
		slog.Error("error creating main world", "error", err)
		os.Exit(1)
	}

	g.SetWorld(mainWorld)

	if err := g.Start(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
