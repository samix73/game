package systems

import (
	"github.com/samix73/game/components"
	"github.com/samix73/game/ecs"
	"github.com/samix73/game/helpers"
	"golang.org/x/image/math/f64"
)

var _ ecs.System = (*Gravity)(nil)

type Gravity struct {
	*ecs.BaseSystem

	dv f64.Vec2
}

func NewGravitySystem(priority int, entityManager *ecs.EntityManager, acceleration f64.Vec2) *Gravity {
	return &Gravity{
		BaseSystem: ecs.NewBaseSystem(ecs.NextID(), priority, entityManager),
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
