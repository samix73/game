package game

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

var _ ebiten.Game = (*Game)(nil)

type Game struct {
	cfg *Config

	activeLevel *Level
}

func NewGame(cfg *Config) *Game {
	return &Game{
		cfg: cfg,
	}
}

func (g *Game) Start() error {
	ebiten.SetWindowSize(g.cfg.ScreenWidth, g.cfg.ScreenHeight)
	ebiten.SetFullscreen(g.cfg.Fullscreen)
	ebiten.SetWindowTitle(g.cfg.Title)

	if err := ebiten.RunGameWithOptions(g, nil); err != nil {
		return fmt.Errorf("game.Game.Start ebiten.RunGameWithOptions error: %w", err)
	}

	return nil
}

func (g *Game) SetActiveLevel(level *Level) {
	g.activeLevel = level
}

func (g *Game) Layout(outsideWidth int, outsideHeight int) (int, int) {
	return g.cfg.ScreenWidth, g.cfg.ScreenHeight
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.activeLevel != nil {
		g.activeLevel.Draw(screen)
	}
}

func (g *Game) Update() error {
	if g.activeLevel != nil {
		if err := g.activeLevel.Update(); err != nil {
			return fmt.Errorf("game.Game.Update activeLevel.Update error: %w", err)
		}
	}

	return nil
}
