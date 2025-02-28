package game

import (
	"fmt"
	"os"
	"runtime/trace"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

var _ ebiten.Game = (*Game)(nil)

type Config struct {
	Title                     string
	ScreenWidth, ScreenHeight int
	Fullscreen                bool
	Tracing                   bool
	TraceFile                 string
}

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

	closer, err := g.setupTrace()
	if err != nil {
		return fmt.Errorf("game.Game.Start g.setupTrace error: %w", err)
	}

	if closer != nil {
		defer closer()
	}

	if err := ebiten.RunGameWithOptions(g, nil); err != nil {
		return fmt.Errorf("game.Game.Start ebiten.RunGameWithOptions error: %w", err)
	}

	return nil
}

func (g *Game) setupTrace() (func(), error) {
	if g.cfg.Tracing {
		filename := fmt.Sprintf("trace_%s.out",
			time.Now().UTC().Format("2006-01-02_15-04-05"),
		)
		f, err := os.Create(filename)
		if err != nil {
			return nil, fmt.Errorf("game.Game.setupTrace os.Create error: %w", err)
		}

		if err := trace.Start(f); err != nil {
			return nil, fmt.Errorf("game.Game.setupTrace trace.Start error: %w", err)
		}

		return func() {
			f.Close()
			trace.Stop()
		}, nil
	}

	return nil, nil
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
