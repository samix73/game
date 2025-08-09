package entities

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/png"
	"runtime/trace"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/samix73/game/assets"
	"github.com/samix73/game/components"
	"github.com/samix73/game/ecs"
)

func NewBiogEntity(ctx context.Context, em *ecs.EntityManager) (ecs.EntityID, error) {
	ctx, task := trace.NewTask(ctx, "entities.NewBiogEntity")
	defer task.End()

	f, err := assets.GetSprite("biog")
	if err != nil {
		return ecs.UndefinedID, fmt.Errorf("error getting sprite: %v", err)
	}

	img, _, err := image.Decode(bytes.NewReader(f))
	if err != nil {
		return ecs.UndefinedID, fmt.Errorf("error decoding image: %v", err)
	}

	entity := em.NewEntity(ctx)
	ecs.AddComponent[components.Transform](ctx, em, entity)

	rigidbody := ecs.AddComponent[components.RigidBody](ctx, em, entity)
	rigidbody.Gravity = false

	renderable := ecs.AddComponent[components.Renderable](ctx, em, entity)
	renderable.Sprite = ebiten.NewImageFromImage(img)

	return entity, nil
}
