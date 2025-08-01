package systems

import (
	"github.com/samix73/game/components"
	"github.com/samix73/game/ecs"
	"golang.org/x/image/math/f64"
)

var _ ecs.System = (*Gravity)(nil)

type Gravity struct {
	*ecs.BaseSystem
	Acceleration f64.Vec2
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

		rigidBody.Velocity[0] += g.Acceleration[0]
		rigidBody.Velocity[1] += g.Acceleration[1]
	}

	return nil
}
