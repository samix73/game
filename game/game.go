package game

import (
	"fmt"
	"os"
	"runtime/trace"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/samix73/game/ecs"
	"golang.org/x/image/math/f64"
)

var _ ebiten.Game = (*Game)(nil)

type Config struct {
	Title                     string
	Gravity                   f64.Vec2
	ScreenWidth, ScreenHeight int
	Fullscreen                bool
	Tracing                   bool
	TraceFile                 string
}

type Game struct {
	cfg *Config

	activeWorld ecs.World
}

func NewGame(cfg *Config) *Game {
	return &Game{
		cfg: cfg,
	}
}

func (g *Game) Config() *Config {
	return g.cfg
}

func (g *Game) SetWorld(world ecs.World) {
	if g.activeWorld != nil {
		g.activeWorld.Teardown()
	}

	g.activeWorld = world
}

func (g *Game) Start() error {
	ebiten.SetWindowSize(g.cfg.ScreenWidth, g.cfg.ScreenHeight)
	ebiten.SetFullscreen(g.cfg.Fullscreen)
	ebiten.SetWindowTitle(g.cfg.Title)

	closer, err := g.setupTrace()
	if err != nil {
		return fmt.Errorf("game.Game.Start g.setupTrace error: %w", err)
	}

	defer closer()

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
			_ = f.Close()
			trace.Stop()
		}, nil
	}

	return nil, nil
}

func (g *Game) Layout(outsideWidth int, outsideHeight int) (int, int) {
	return g.cfg.ScreenWidth, g.cfg.ScreenHeight
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.activeWorld == nil {
		return
	}

	g.activeWorld.Draw(screen)
}

func (g *Game) Update() error {
	if g.activeWorld == nil {
		return nil
	}

	if err := g.activeWorld.Update(); err != nil {
		return fmt.Errorf("game.Game.Update activeWorld.Update error: %w", err)
	}

	return nil
}
