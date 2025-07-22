package worlds

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/samix73/game/ecs"
	"github.com/samix73/game/internal/game"
)

var _ ecs.World = (*MainWorld)(nil)

type MainWorld struct {
	*ecs.BaseWorld

	g *game.Game
}

func NewMainWorld(g *game.Game) *MainWorld {
	entityManager := ecs.NewEntityManager()
	systemManager := ecs.NewSystemManager(entityManager)

	return &MainWorld{
		BaseWorld: ecs.NewBaseWorld(entityManager, systemManager),
		g:         g,
	}
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
