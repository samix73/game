package worlds

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/samix73/game/ecs"
	"github.com/samix73/game/game"
	"github.com/samix73/game/systems"
)

var _ ecs.World = (*MainWorld)(nil)

type MainWorld struct {
	*ecs.BaseWorld

	g *game.Game
}

func NewMainWorld(g *game.Game) *MainWorld {
	entityManager := ecs.NewEntityManager()
	systemManager := ecs.NewSystemManager(entityManager)

	w := &MainWorld{
		BaseWorld: ecs.NewBaseWorld(entityManager, systemManager),
		g:         g,
	}

	w.registerSystems()

	return w
}

func (m *MainWorld) registerSystems() {
	m.SystemManager().Add(
		systems.NewCameraSystem(0, m.EntityManager()),
	)
}

func (m *MainWorld) Draw(screen *ebiten.Image) {
	m.SystemManager().Draw(screen)
}

func (m *MainWorld) Update() error {
	if err := m.SystemManager().Update(); err != nil {
		return fmt.Errorf("error updating systems: %w", err)
	}

	return nil
}

func (m *MainWorld) Teardown() {
	m.SystemManager().Teardown()
}
