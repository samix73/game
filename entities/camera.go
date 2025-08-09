package entities

import (
	"context"
	"runtime/trace"

	"github.com/samix73/game/components"
	"github.com/samix73/game/ecs"
)

func NewCameraEntity(ctx context.Context, em *ecs.EntityManager, active bool) ecs.EntityID {
	ctx, task := trace.NewTask(ctx, "entities.NewCameraEntity")
	defer task.End()

	entity := em.NewEntity(ctx)
	if active {
		ecs.AddComponent[components.ActiveCamera](ctx, em, entity)
	}
	ecs.AddComponent[components.Transform](ctx, em, entity)
	camera := ecs.AddComponent[components.Camera](ctx, em, entity)
	camera.Zoom = 1.0

	return entity
}
