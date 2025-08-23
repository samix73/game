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

func (m *MainWorld) Teardown() {
	m.SystemManager().Teardown()
}
