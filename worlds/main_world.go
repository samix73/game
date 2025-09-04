package worlds

import (
	ecs "github.com/samix73/ebiten-ecs"
	"github.com/samix73/game/systems"
)

var _ ecs.World = (*MainWorld)(nil)

type MainWorld struct {
	*ecs.BaseWorld
}

func (m *MainWorld) Init(g *ecs.Game) error {
	g.SetTimeScale(1)

	entityManager := ecs.NewEntityManager()
	systemManager := ecs.NewSystemManager(entityManager, g)

	m.BaseWorld = ecs.NewBaseWorld(systemManager)

	m.registerSystems()

	return nil
}

func (m *MainWorld) registerSystems() {
	m.SystemManager().Add(
		systems.NewPauseSystem(0),
		systems.NewGravitySystem(1),
		systems.NewPhysicsSystem(2),
		systems.NewCollisionSystem(3),
		systems.NewCameraSystem(4),
	)
}
