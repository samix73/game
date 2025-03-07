package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/samix73/game/internal/game"
)

var (
	fullscreen = flag.Bool("fullscreen", false, "enable fullscreen mode")
	tracing    = flag.Bool("tracing", false, "enable tracing")
)

func main() {
	flag.Parse()

	g := game.NewGame(&game.Config{
		Title:        "Game",
		ScreenWidth:  1280,
		ScreenHeight: 960,
		Fullscreen:   *fullscreen,
		Tracing:      *tracing,
	})

	if err := g.Start(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
