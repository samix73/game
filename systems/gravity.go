package systems

import (
	ecs "github.com/samix73/ebiten-ecs"
	"github.com/samix73/game/components"
	"golang.org/x/image/math/f64"
)

var gravity = f64.Vec2{0, -981}

var _ ecs.System = (*Gravity)(nil)

type Gravity struct {
	*ecs.BaseSystem

	dv f64.Vec2
}

func NewGravitySystem(priority int) *Gravity {
	return &Gravity{
		BaseSystem: ecs.NewBaseSystem(ecs.NextID(), priority),
		dv:         gravity,
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
