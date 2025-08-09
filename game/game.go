package game

import (
	"context"
	"fmt"
	"runtime/trace"

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

	PlayerJumpKey             ebiten.Key
	PlayerJumpForce           float64
	PlayerForwardAcceleration float64
}

type Game struct {
	cfg *Config
	ctx context.Context

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

func (g *Game) Layout(outsideWidth int, outsideHeight int) (int, int) {
	return g.cfg.ScreenWidth, g.cfg.ScreenHeight
}

func (g *Game) Draw(screen *ebiten.Image) {
	ctx, task := trace.NewTask(g.ctx, "game.Game.Draw")
	defer task.End()

	if g.activeWorld == nil {
		return
	}

	g.activeWorld.Draw(ctx, screen)
}

func (g *Game) Update() error {
	ctx, task := trace.NewTask(g.ctx, "game.Game.Update")
	defer task.End()

	if g.activeWorld == nil {
		return nil
	}

	if err := g.activeWorld.Update(ctx); err != nil {
		return fmt.Errorf("game.Game.Update activeWorld.Update error: %w", err)
	}

	return nil
}
