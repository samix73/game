package game

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

var _ ebiten.Game = (*Game)(nil)

type Game struct {
	cfg *Config
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

func (g *Game) Layout(outsideWidth int, outsideHeight int) (int, int) {
	return g.cfg.ScreenWidth, g.cfg.ScreenHeight
}

func (g *Game) Draw(screen *ebiten.Image) {

}

func (g *Game) Update() error {
	return nil
}
