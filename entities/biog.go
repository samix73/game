package entities

import (
	"context"
	"fmt"
	_ "image/png"
	"runtime/trace"

	"github.com/samix73/game/assets"
	"github.com/samix73/game/components"
	"github.com/samix73/game/ecs"
)

func NewBiogEntity(ctx context.Context, em *ecs.EntityManager) (ecs.EntityID, error) {
	ctx, task := trace.NewTask(ctx, "entities.NewBiogEntity")
	defer task.End()

	img, err := assets.GetSprite(ctx, "biog.png")
	if err != nil {
		return ecs.UndefinedID, fmt.Errorf("error getting sprite: %v", err)
	}

	entity := em.NewEntity(ctx)
	ecs.AddComponent[components.Transform](ctx, em, entity)
	ecs.AddComponent[components.Player](ctx, em, entity)
	collider := ecs.AddComponent[components.ColliderComponent](ctx, em, entity)
	collider.Bounds.SetImageBounds(img.Bounds())

	rigidBody := ecs.AddComponent[components.RigidBody](ctx, em, entity)
	rigidBody.Gravity = true

	renderable := ecs.AddComponent[components.Renderable](ctx, em, entity)
	renderable.Sprite = img

	return entity, nil
}
