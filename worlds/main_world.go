package worlds

import (
	"fmt"
	"math/rand/v2"

	"github.com/samix73/game/ecs"
	"github.com/samix73/game/entities"
	"github.com/samix73/game/game"
	"github.com/samix73/game/systems"
	"golang.org/x/image/math/f64"
)

var _ ecs.World = (*MainWorld)(nil)

type MainWorld struct {
	*ecs.BaseWorld

	g *game.Game
}

func NewMainWorld(g *game.Game) (*MainWorld, error) {
	entityManager := ecs.NewEntityManager()
	systemManager := ecs.NewSystemManager(entityManager)

	if _, err := entities.NewBiogEntity(entityManager); err != nil {
		return nil, fmt.Errorf("error creating biog entity: %w", err)
	}

	colors := []string{"red", "yellow", "blue"}
	for i := range 10_000 {
		if _, err := entities.NewObstacleEntity(
			entityManager,
			colors[rand.IntN(len(colors))],
			rand.IntN(8)+3,
			f64.Vec2{float64(i * 450), 0},
		); err != nil {
			return nil, fmt.Errorf("error creating obstacle entity: %w", err)
		}
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
		systems.NewPhysicsSystem(1, m.EntityManager()),
		systems.NewGravitySystem(2, m.EntityManager(), gameCfg.Gravity),
		systems.NewCollisionSystem(3, m.EntityManager()),
		systems.NewPlayerSystem(4, m.EntityManager(),
			gameCfg.PlayerJumpForce, gameCfg.PlayerForwardAcceleration, gameCfg.PlayerCameraOffset, gameCfg.PlayerMaxSpeed),
	)
}

func (m *MainWorld) Teardown() {
	m.SystemManager().Teardown()
}
