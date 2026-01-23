package ecs

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/samix73/game/game/assets"
)

var _ ebiten.Game = (*Game)(nil)

type GameConfig struct {
	Title                     string
	ScreenWidth, ScreenHeight int
	Fullscreen                bool
}

type Game struct {
	cfg         *GameConfig
	activeWorld *World
	timeScale   float64
}

func NewGame(cfg *GameConfig) *Game {
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

func (g *Game) Config() GameConfig {
	return *g.cfg
}

func (g *Game) SetActiveWorld(world *World) error {
	if g.activeWorld != nil {
		g.activeWorld.Teardown()
	}

	g.activeWorld = world

	return nil
}

func (g *Game) loadSystems(systemManager *SystemManager, systemCfgs []SystemConfig) error {
	for _, systemCfg := range systemCfgs {
		systemCtor, ok := GetSystem(systemCfg.Name)
		if !ok {
			return fmt.Errorf("ecs.LoadWorld: system %s not found", systemCfg.Name)
		}

		systemManager.Add(systemCtor(systemCfg.Priority))
	}

	return nil
}

// loadEntities loads entities from entity configs and world metadata.
// It overwrites components from entity configs with components from world metadata.
func (g *Game) loadEntities(em *EntityManager, entityCfgs []EntityConfig, worldMD toml.MetaData) error {
	for _, entityCfg := range entityCfgs {
		entity := em.NewEntity()

		if entityCfg.Name == "" {
			return errors.New("ecs.LoadWorld: entity name is empty")
		}

		entityData, err := assets.GetEntity("asd")
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("ecs.LoadWorld: %w", err)
		}

		// Load components from entity config
		if len(entityData) > 0 {
			components := make(EntityComponentsConfig)
			md, err := toml.NewDecoder(bytes.NewReader(entityData)).Decode(&components)
			if err != nil {
				return fmt.Errorf("ecs.LoadWorld: %w", err)
			}

			for name, args := range components {
				if name == "" {
					return errors.New("ecs.LoadWorld: component name is empty")
				}

				component, ok := NewComponent(em, name)
				if !ok {
					return fmt.Errorf("ecs.LoadWorld: component %s not found", name)
				}

				if err := md.PrimitiveDecode(args, component); err != nil {
					return fmt.Errorf("ecs.LoadWorld: failed to decode component %s: %w",
						name, err)
				}

				if err := em.AddComponent(entity, component); err != nil {
					return fmt.Errorf("ecs.LoadWorld: failed to add component %s to entity %s: %w",
						name, entityCfg.Name, err)
				}
			}
		}

		// Overwrite components from entity config
		for name, args := range entityCfg.Components {
			if name == "" {
				return errors.New("ecs.LoadWorld: component name is empty")
			}

			component, ok := NewComponent(em, name)
			if !ok {
				return fmt.Errorf("ecs.LoadWorld: component %s not found", name)
			}

			if err := worldMD.PrimitiveDecode(args, component); err != nil {
				return fmt.Errorf("ecs.LoadWorld: failed to decode component %s: %w",
					name, err)
			}

			if err := em.AddComponent(entity, component); err != nil {
				return fmt.Errorf("ecs.LoadWorld: failed to add component %s to entity %s: %w",
					name, entityCfg.Name, err)
			}
		}
	}

	return nil
}

func (g *Game) LoadWorld(path string) (*World, error) {
	data, err := assets.GetWorld(path)
	if err != nil {
		return nil, fmt.Errorf("ecs.LoadWorld: %w", err)
	}

	var worldConfig WorldConfig

	md, err := toml.NewDecoder(bytes.NewReader(data)).Decode(&worldConfig)
	if err != nil {
		return nil, fmt.Errorf("ecs.LoadWorld: %w", err)
	}

	em := NewEntityManager()
	sm := NewSystemManager(em, g)

	if err := g.loadSystems(sm, worldConfig.Systems); err != nil {
		return nil, fmt.Errorf("ecs.LoadWorld: %w", err)
	}

	if err := g.loadEntities(em, worldConfig.Entities, md); err != nil {
		return nil, fmt.Errorf("ecs.LoadWorld: %w", err)
	}

	return &World{
		cfg: worldConfig,

		systemManager: sm,
		entityManager: em,
	}, nil
}

func (g *Game) DeltaTime() float64 {
	return 1.0 / float64(ebiten.TPS()) * g.TimeScale()
}

func (g *Game) Start() error {
	ebiten.SetWindowSize(g.cfg.ScreenWidth, g.cfg.ScreenHeight)
	ebiten.SetFullscreen(g.cfg.Fullscreen)
	ebiten.SetWindowTitle(g.cfg.Title)

	if err := ebiten.RunGameWithOptions(g, nil); err != nil {
		return fmt.Errorf("ecs.Game.Start ebiten.RunGameWithOptions error: %w", err)
	}

	return nil
}

func (g *Game) Layout(outsideWidth int, outsideHeight int) (screenWidth int, screenHeight int) {
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
		return fmt.Errorf("ecs.Game.Update activeWorld.Update error: %w", err)
	}

	return nil
}
