package game

import (
	"context"
	"fmt"
	"runtime/trace"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/samix73/game/ecs"
	"github.com/samix73/game/keys"
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
	cfg    *Config
	ctx    context.Context
	paused bool

	activeWorld ecs.World
}

func NewGame(ctx context.Context, cfg *Config) *Game {
	return &Game{
		cfg: cfg,
		ctx: ctx,
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

	if err := ebiten.RunGameWithOptions(g, nil); err != nil {
		return fmt.Errorf("game.Game.Start ebiten.RunGameWithOptions error: %w", err)
	}

	return nil
}

func (g *Game) Pause() bool {
	if keys.IsPressed(keys.PauseAction) {
		g.paused = !g.paused
	}

	return g.paused
}

func (g *Game) Layout(outsideWidth int, outsideHeight int) (int, int) {
	return g.cfg.ScreenWidth, g.cfg.ScreenHeight
}

func (g *Game) Draw(screen *ebiten.Image) {
	ctx, task := trace.NewTask(g.ctx, "game.Game.Draw")
	defer task.End()

	if g.paused {
		ebitenutil.DebugPrintAt(screen, "Paused - press P to resume", 16, 16)
	}

	if g.activeWorld == nil {
		return
	}

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("FPS: %.2f", ebiten.ActualFPS()), 16, 32)

	g.activeWorld.Draw(ctx, screen)
}

func (g *Game) Update() error {
	ctx, task := trace.NewTask(g.ctx, "game.Game.Update")
	defer task.End()

	if g.Pause() {
		return nil
	}

	if g.activeWorld == nil {
		return nil
	}

	if err := g.activeWorld.Update(ctx); err != nil {
		return fmt.Errorf("game.Game.Update activeWorld.Update error: %w", err)
	}

	return nil
}
