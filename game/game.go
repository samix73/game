package game

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/samix73/game/ecs"
	"golang.org/x/image/math/f64"
)

var _ ebiten.Game = (*Game)(nil)

type Config struct {
	Title                     string
	Gravity                   f64.Vec2
	ScreenWidth, ScreenHeight int
	Fullscreen                bool

	PlayerJumpForce           float64
	PlayerForwardAcceleration float64
	PlayerCameraOffset        f64.Vec2
	PlayerMaxSpeed            float64
}

type Game struct {
	cfg *Config

	activeWorld ecs.World
	timeScale   float64
}

func NewGame(cfg *Config) *Game {
	return &Game{
		cfg:       cfg,
		timeScale: 1.0,
	}
}

func (g *Game) TimeScale() float64 {
	return g.timeScale
}

func (g *Game) SetTimeScale(scale float64) {
	g.timeScale = math.Max(scale, 0)
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

func (g *Game) DeltaTime() float64 {
	return 1.0 / float64(ebiten.TPS()) * g.TimeScale()
}

func (g *Game) Restart() {
	// TODO re-initliaze the active world
	// g.SetWorld(nil)
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
	if g.activeWorld == nil {
		return
	}

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("FPS: %.2f", ebiten.ActualFPS()), 16, 32)

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
