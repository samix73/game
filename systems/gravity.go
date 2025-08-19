package systems

import (
	"github.com/samix73/game/components"
	"github.com/samix73/game/ecs"
	"github.com/samix73/game/game"
	"github.com/samix73/game/helpers"
	"golang.org/x/image/math/f64"
)

var _ ecs.System = (*Gravity)(nil)

type Gravity struct {
	*ecs.BaseSystem[*game.Game]

	dv f64.Vec2
}

func NewGravitySystem(priority int, entityManager *ecs.EntityManager, game *game.Game) *Gravity {
	cfg := game.Config()

	acceleration := cfg.Gravity

	return &Gravity{
		BaseSystem: ecs.NewBaseSystem(ecs.NextID(), priority, entityManager, game),
		dv: f64.Vec2{
			acceleration[0] * helpers.DeltaTime,
			acceleration[1] * helpers.DeltaTime,
		},
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

		rigidBody.ApplyAcceleration(g.dv)
	}

	return nil
}
