package main

import (
	"log/slog"
	"os"

	"github.com/samix73/game/internal/game"
)

func main() {
	g := game.NewGame(&game.Config{
		Title:        "Game",
		ScreenWidth:  1280,
		ScreenHeight: 960,
	})

	if err := g.Start(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
