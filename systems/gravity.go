package systems

import (
	"github.com/samix73/game/components"
	"github.com/samix73/game/ecs"
	"github.com/samix73/game/game"
	"golang.org/x/image/math/f64"
)

var _ ecs.System = (*Gravity)(nil)

type Gravity struct {
	*ecs.BaseSystem[*game.Game]

	dv f64.Vec2
}

func NewGravitySystem(priority int, entityManager *ecs.EntityManager, game *game.Game) *Gravity {
	return &Gravity{
		BaseSystem: ecs.NewBaseSystem(ecs.NextID(), priority, entityManager, game),
		dv:         game.Config().Gravity,
	}
}

func (g *Gravity) Teardown() {}

func (g *Gravity) Update() error {
	em := g.EntityManager()
	for entity := range ecs.Query[components.RigidBody](em) {
		rigidBody := ecs.MustGetComponent[components.RigidBody](em, entity)
		if rigidBody == nil {
			continue
		}

		if !rigidBody.Gravity {
			continue
		}

		game := g.Game()

		acc := f64.Vec2{
			g.dv[0] * game.DeltaTime(),
			g.dv[1] * game.DeltaTime(),
		}

		rigidBody.ApplyAcceleration(acc)
	}

	return nil
}
