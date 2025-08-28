package entities

import (
	"fmt"
	"path"

	"github.com/hajimehoshi/ebiten/v2"
	ecs "github.com/samix73/ebiten-ecs"
	"github.com/samix73/game/assets"
	"github.com/samix73/game/components"
	"golang.org/x/image/math/f64"
)

func NewObstacleEntity(em *ecs.EntityManager,
	color string, height int, position f64.Vec2) (ecs.EntityID, error) {

	tileName := path.Join("Tiles", "Default", fmt.Sprintf("block_%s.png", color))
	tile, err := assets.GetSprite(tileName)
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

	entity := em.NewEntity()

	transform := ecs.AddComponent[components.Transform](em, entity)
	transform.SetPosition(position[0], position[1])

	renderable := ecs.AddComponent[components.Renderable](em, entity)
	renderable.Sprite = colImg

	obstacle := ecs.AddComponent[components.Obstacle](em, entity)
	obstacle.Color = color
	obstacle.Height = height

	collider := ecs.AddComponent[components.Collider](em, entity)
	collider.Bounds.SetSize(float64(colImg.Bounds().Dx()), float64(colImg.Bounds().Dy()))

	return entity, nil
}
