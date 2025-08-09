package worlds

import (
	"context"
	"fmt"
	"runtime/trace"

	"github.com/samix73/game/ecs"
	"github.com/samix73/game/entities"
	"github.com/samix73/game/game"
	"github.com/samix73/game/systems"
)

var _ ecs.World = (*MainWorld)(nil)

type MainWorld struct {
	*ecs.BaseWorld

	g *game.Game
}

func NewMainWorld(ctx context.Context, g *game.Game) (*MainWorld, error) {
	ctx, task := trace.NewTask(ctx, "worlds.NewMainWorld")
	defer task.End()

	entityManager := ecs.NewEntityManager()
	systemManager := ecs.NewSystemManager(ctx, entityManager)

	if _, err := entities.NewBiogEntity(ctx, entityManager); err != nil {
		return nil, fmt.Errorf("error creating biog entity: %w", err)
	}

	w := &MainWorld{
		BaseWorld: ecs.NewBaseWorld(entityManager, systemManager),
		g:         g,
	}

	w.registerSystems(ctx)

	return w, nil
}

func (m *MainWorld) registerSystems(ctx context.Context) {
	ctx, task := trace.NewTask(ctx, "worlds.MainWorld.registerSystems")
	defer task.End()

	gameCfg := m.g.Config()
	m.SystemManager().Add(ctx,
		systems.NewCameraSystem(ctx, 0, m.EntityManager(), gameCfg.ScreenWidth, gameCfg.ScreenHeight),
		systems.NewPhysicsSystem(ctx, 1, m.EntityManager()),
		systems.NewGravitySystem(ctx, 2, m.EntityManager(), gameCfg.Gravity),
		systems.NewPlayerSystem(ctx, 3, m.EntityManager(),
			gameCfg.PlayerJumpKey, gameCfg.PlayerJumpForce, gameCfg.PlayerForwardAcceleration),
	)
}

func (m *MainWorld) Teardown() {
	m.SystemManager().Teardown()
}
