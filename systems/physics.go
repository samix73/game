package systems

import (
	"context"
	"runtime/trace"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/samix73/game/components"
	"github.com/samix73/game/ecs"
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

	deltaTime := 1.0 / ebiten.ActualTPS()

	em := p.EntityManager()

	for entity := range ecs.Query2[components.RigidBody, components.Transform](ctx, em) {
		rigidBody := ecs.MustGetComponent[components.RigidBody](ctx, em, entity)
		transform := ecs.MustGetComponent[components.Transform](ctx, em, entity)

		newX := transform.Vec2[0] + rigidBody.Velocity[0]*deltaTime
		newY := transform.Vec2[1] + rigidBody.Velocity[1]*deltaTime

		transform.Vec2[0] = newX
		transform.Vec2[1] = newY
	}

	return nil
}
