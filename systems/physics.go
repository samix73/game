package systems

import (
	"context"
	"log/slog"
	"runtime/trace"

	"github.com/samix73/game/components"
	"github.com/samix73/game/ecs"
	"github.com/samix73/game/helpers"
	"golang.org/x/image/math/f64"
)

type Physics struct {
	*ecs.BaseSystem
}

func NewPhysicsSystem(ctx context.Context, priority int, entityManager *ecs.EntityManager) *Physics {
	ctx, task := trace.NewTask(ctx, "systems.NewPhysicsSystem")
	defer task.End()

	return &Physics{
		BaseSystem: ecs.NewBaseSystem(ctx, ecs.NextID(ctx), priority, entityManager),
	}
}

func (p *Physics) Teardown() {}

func (p *Physics) Update(ctx context.Context) error {
	ctx, task := trace.NewTask(ctx, "systems.Physics.Update")
	defer task.End()

	em := p.EntityManager()

	for entity := range ecs.Query2[components.RigidBody, components.Transform](ctx, em) {
		rigidBody := ecs.MustGetComponent[components.RigidBody](ctx, em, entity)
		transform := ecs.MustGetComponent[components.Transform](ctx, em, entity)

		transform.Translate(f64.Vec2{
			rigidBody.Velocity[0] * helpers.DeltaTime,
			rigidBody.Velocity[1] * helpers.DeltaTime,
		})

		slog.Debug("Physics.Update",
			slog.Uint64("entity", uint64(entity)),
			slog.Any("position", transform.Vec2),
			slog.Any("velocity", rigidBody.Velocity),
		)
	}

	return nil
}
