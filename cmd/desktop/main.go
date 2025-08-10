package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"runtime/trace"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/samix73/game/game"
	"github.com/samix73/game/worlds"
	"golang.org/x/image/math/f64"
)

var (
	fullscreen   = flag.Bool("fullscreen", false, "enable fullscreen mode")
	traceEnabled = flag.Bool("trace", false, "enable tracing")
	logLevel     = flag.String("log-level", "info", "set the log level (debug, info, warn, error, fatal)")
)

func setupTrace(enabled bool) (func(), error) {
	if !enabled {
		return nil, nil
	}

	filename := fmt.Sprintf("trace_%s.out",
		time.Now().Format("2006-01-02_15-04-05"),
	)
	f, err := os.Create(filename)
	if err != nil {
		return nil, fmt.Errorf("game.Game.setupTrace os.Create error: %w", err)
	}

	if err := trace.Start(f); err != nil {
		return nil, fmt.Errorf("game.Game.setupTrace trace.Start error: %w", err)
	}

	return func() {
		_ = f.Close()
		trace.Stop()
	}, nil
}

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

	ctx, cancel := context.WithCancel(context.Background())

	closeTrace, err := setupTrace(*traceEnabled)
	if err != nil {
		slog.Error("error setting up trace", slog.String("error", err.Error()))
		os.Exit(1)
	}

	defer func() {
		cancel()
		if closeTrace != nil {
			closeTrace()
		}
	}()

	g := game.NewGame(ctx, &game.Config{
		Title:        "Game",
		ScreenWidth:  1280,
		ScreenHeight: 960,
		Gravity:      f64.Vec2{0, 981},
		Fullscreen:   *fullscreen,

		PlayerJumpKey:             ebiten.KeySpace,
		PlayerJumpForce:           500,
		PlayerForwardAcceleration: 35,
		PlayerCameraOffset:        f64.Vec2{200, 0},
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
