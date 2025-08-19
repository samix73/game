package worlds

import (
	"fmt"

	"github.com/samix73/game/ecs"
	"github.com/samix73/game/entities"
	"github.com/samix73/game/game"
	"github.com/samix73/game/systems"
)

var _ ecs.World = (*MainWorld)(nil)

type MainWorld struct {
	*ecs.BaseWorld[*game.Game]
}

func NewMainWorld(g *game.Game) (*MainWorld, error) {
	entityManager := ecs.NewEntityManager()
	systemManager := ecs.NewSystemManager(entityManager)

	if _, err := entities.NewBiogEntity(entityManager); err != nil {
		return nil, fmt.Errorf("error creating biog entity: %w", err)
	}

	w := &MainWorld{
		BaseWorld: ecs.NewBaseWorld(entityManager, systemManager, g),
	}

	w.registerSystems()

	return w, nil
}

func (m *MainWorld) registerSystems() {
	gameCfg := m.Game().Config()
	m.SystemManager().Add(
		systems.NewPlayerSystem(0, m.EntityManager(), m.Game(),
			gameCfg.PlayerJumpForce, gameCfg.PlayerForwardAcceleration, gameCfg.PlayerCameraOffset, gameCfg.PlayerMaxSpeed),
		systems.NewGravitySystem(1, m.EntityManager(), m.Game(), gameCfg.Gravity),
		systems.NewPhysicsSystem(2, m.EntityManager(), m.Game()),
		systems.NewCollisionSystem(3, m.EntityManager(), m.Game()),
		systems.NewPlayerCollisionSystem(4, m.EntityManager(), m.Game()),
		systems.NewLevelGenSystem(5, m.EntityManager(), m.Game()),
		systems.NewCameraSystem(6, m.EntityManager(), m.Game(), gameCfg.ScreenWidth, gameCfg.ScreenHeight),
	)
}

func (m *MainWorld) Teardown() {
	m.SystemManager().Teardown()
}
