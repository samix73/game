package physics

import (
	"github.com/jakecoffman/cp"
	ecs "github.com/samix73/ebiten-ecs"
	"github.com/samix73/game/client/components"
)

var gravity = cp.Vector{X: 0, Y: -981}

var _ ecs.System = (*Gravity)(nil)

type Gravity struct {
	*ecs.BaseSystem

	dv cp.Vector
}

func NewGravitySystem(priority int) *Gravity {
	return &Gravity{
		BaseSystem: ecs.NewBaseSystem(ecs.NextID(), priority),
		dv:         gravity,
	}
}

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

		acc := cp.Vector{
			X: g.dv.X * game.DeltaTime(),
			Y: g.dv.Y * game.DeltaTime(),
		}

		rigidBody.ApplyAcceleration(acc)
	}

	return nil
}
