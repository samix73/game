package main

import (
	"flag"
	"log/slog"
	"os"

	ecs "github.com/samix73/ebiten-ecs"
	"github.com/samix73/game/worlds"
)

var (
	fullscreen = flag.Bool("fullscreen", false, "enable fullscreen mode")
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

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))
	slog.SetDefault(logger)
}

func main() {
	flag.Parse()

	setupLogger(*logLevel)

	g := ecs.NewGame(&ecs.GameConfig{
		Title:        "Game",
		ScreenWidth:  1280,
		ScreenHeight: 960,
		Fullscreen:   *fullscreen,
	})

	if err := g.SetActiveWorld(new(worlds.MainWorld)); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	if err := g.Start(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
