package systems

import (
	"context"
	"runtime/trace"

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

func NewGravitySystem(ctx context.Context, priority int, entityManager *ecs.EntityManager, acceleration f64.Vec2) *Gravity {
	ctx, task := trace.NewTask(ctx, "systems.NewGravitySystem")
	defer task.End()

	return &Gravity{
		BaseSystem: ecs.NewBaseSystem(ctx, ecs.NextID(ctx), priority, entityManager),
		dv: f64.Vec2{
			acceleration[0] * helpers.DeltaTime,
			acceleration[1] * helpers.DeltaTime,
		},
	}
}

func (g *Gravity) Teardown() {}

func (g *Gravity) Update(ctx context.Context) error {
	ctx, task := trace.NewTask(ctx, "systems.Gravity.Update")
	defer task.End()

	em := g.EntityManager()
	for entity := range ecs.Query[components.RigidBody](ctx, em) {
		rigidBody := ecs.MustGetComponent[components.RigidBody](ctx, em, entity)
		if rigidBody == nil {
			continue
		}

		if !rigidBody.Gravity {
			continue
		}

		rigidBody.Velocity[0] += g.dv[0]
		rigidBody.Velocity[1] += g.dv[1]
	}

	return nil
}
