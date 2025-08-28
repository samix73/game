package worlds

import (
	"fmt"

	"github.com/samix73/ebiten-ecs"
	"github.com/samix73/game/entities"
	"github.com/samix73/game/systems"
)

var _ ecs.World = (*MainWorld)(nil)

type MainWorld struct {
	*ecs.BaseWorld
}

func (m *MainWorld) Init(g *ecs.Game) error {
	g.SetTimeScale(1)

	entityManager := ecs.NewEntityManager()
	systemManager := ecs.NewSystemManager(entityManager)

	m.BaseWorld = ecs.NewBaseWorld(entityManager, systemManager, g)

	if _, err := entities.NewBiogEntity(entityManager); err != nil {
		return fmt.Errorf("error creating biog entity: %w", err)
	}

	m.registerSystems()

	return nil
}

func (m *MainWorld) registerSystems() {
	m.SystemManager().Add(
		systems.NewPauseSystem(0, m.EntityManager(), m.Game()),
		systems.NewPlayerSystem(1, m.EntityManager(), m.Game()),
		systems.NewGravitySystem(2, m.EntityManager(), m.Game()),
		systems.NewPhysicsSystem(3, m.EntityManager(), m.Game()),
		systems.NewCollisionSystem(4, m.EntityManager(), m.Game()),
		systems.NewPlayerCollisionSystem(5, m.EntityManager(), m.Game()),
		systems.NewLevelGenSystem(6, m.EntityManager(), m.Game()),
		systems.NewCameraSystem(7, m.EntityManager(), m.Game()),
		systems.NewRestartSystem(8, m.EntityManager(), m.Game()),
	)
}
