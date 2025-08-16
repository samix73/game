package entities

import (
	"context"
	"fmt"
	"path"
	"runtime/trace"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/samix73/game/assets"
	"github.com/samix73/game/components"
	"github.com/samix73/game/ecs"
	"golang.org/x/image/math/f64"
)

func NewObstacleEntity(ctx context.Context, em *ecs.EntityManager,
	color string, height int, position f64.Vec2) (ecs.EntityID, error) {
	ctx, task := trace.NewTask(ctx, "entities.NewObstacleEntity")
	defer task.End()

	tileName := path.Join("Tiles", "Default", fmt.Sprintf("block_%s.png", color))
	tile, err := assets.GetSprite(ctx, tileName)
	if err != nil {
		return 0, fmt.Errorf("entities.NewObstacleEntity GetSprite error: %w", err)
	}

	tw := tile.Bounds().Dx()
	th := tile.Bounds().Dy()
	if height < 1 {
		height = 1
	}

	colImg := ebiten.NewImage(tw, th*height)
	for i := range height {
		var op ebiten.DrawImageOptions
		op.GeoM.Translate(0, float64(i*th))
		colImg.DrawImage(tile, &op)
	}

	entity := em.NewEntity(ctx)

	transform := ecs.AddComponent[components.Transform](ctx, em, entity)
	transform.SetPosition(position[0], position[1])

	renderable := ecs.AddComponent[components.Renderable](ctx, em, entity)
	renderable.Sprite = colImg

	obstacle := ecs.AddComponent[components.Obstacle](ctx, em, entity)
	obstacle.Color = color
	obstacle.Height = height

	collider := ecs.AddComponent[components.ColliderComponent](ctx, em, entity)
	collider.Bounds.SetImageBounds(colImg.Bounds())

	return entity, nil
}
