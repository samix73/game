package worlds

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
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

func NewMainWorld(g *game.Game) (*MainWorld, error) {
	entityManager := ecs.NewEntityManager()
	systemManager := ecs.NewSystemManager(entityManager)

	if _, err := entities.NewDrawMeEntity(entityManager); err != nil {
		return nil, fmt.Errorf("error creating draw me entity: %w", err)
	}

	w := &MainWorld{
		BaseWorld: ecs.NewBaseWorld(entityManager, systemManager),
		g:         g,
	}

	w.registerSystems()

	return w, nil
}

func (m *MainWorld) registerSystems() {
	gameCfg := m.g.Config()
	m.SystemManager().Add(
		systems.NewCameraSystem(0, m.EntityManager(), gameCfg.ScreenWidth, gameCfg.ScreenHeight),
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
