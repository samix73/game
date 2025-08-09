package systems

import (
	"context"
	"fmt"
	"log/slog"
	"runtime/trace"

	"github.com/samix73/game/components"
	"github.com/samix73/game/ecs"
	"github.com/samix73/game/helpers"
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

		transform.Vec2[0] += rigidBody.Velocity[0] * helpers.DeltaTime
		transform.Vec2[1] += rigidBody.Velocity[1] * helpers.DeltaTime

		slog.Debug("Physics.Update",
			slog.Uint64("entity", uint64(entity)),
			slog.String("position", fmt.Sprintf("(%.2f, %.2f)", transform.Vec2[0], transform.Vec2[1])),
			slog.String("velocity", fmt.Sprintf("(%.2f, %.2f)", rigidBody.Velocity[0], rigidBody.Velocity[1])),
		)
	}

	return nil
}
