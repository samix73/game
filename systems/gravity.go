package systems

import (
	"context"
	"runtime/trace"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/samix73/game/components"
	"github.com/samix73/game/ecs"
	"golang.org/x/image/math/f64"
)

var _ ecs.System = (*Gravity)(nil)

type Gravity struct {
	*ecs.BaseSystem
	Acceleration f64.Vec2
}

func NewGravitySystem(ctx context.Context, priority int, entityManager *ecs.EntityManager, acceleration f64.Vec2) *Gravity {
	ctx, task := trace.NewTask(ctx, "systems.NewGravitySystem")
	defer task.End()

	return &Gravity{
		BaseSystem:   ecs.NewBaseSystem(ctx, ecs.NextID(ctx), priority, entityManager),
		Acceleration: acceleration,
	}
}

func (g *Gravity) Teardown() {}

func (g *Gravity) Update(ctx context.Context) error {
	ctx, task := trace.NewTask(ctx, "systems.Gravity.Update")
	defer task.End()

	deltaTime := 1.0 / ebiten.ActualTPS()

	em := g.EntityManager()
	for entity := range ecs.Query[components.RigidBody](ctx, em) {
		rigidBody := ecs.MustGetComponent[components.RigidBody](ctx, em, entity)
		if rigidBody == nil {
			continue
		}

		if !rigidBody.Gravity {
			continue
		}

		rigidBody.Velocity[0] += g.Acceleration[0] * deltaTime
		rigidBody.Velocity[1] += g.Acceleration[1] * deltaTime
	}

	return nil
}
