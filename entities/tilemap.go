package entities

import (
	"bytes"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/images"
	ecs "github.com/samix73/ebiten-ecs"
	"github.com/samix73/game/components"
)

const (
	tileSize = 16
)

func NewTileMapEntity(em *ecs.EntityManager) ecs.EntityID {
	entity := em.NewEntity()

	ecs.AddComponent[components.Transform](em, entity)
	ecs.AddComponent[components.Renderable](em, entity)
	tileMap := ecs.AddComponent[components.TileMap](em, entity)

	img, _, _ := image.Decode(bytes.NewReader(images.Tiles_png))
	atlas := ebiten.NewImageFromImage(img)
	tileMap.SetAtlas(atlas)
	tileMap.TileSize = tileSize
	tileMap.Layer = 0
	tileMap.Width = 240
	tileMap.Height = 240

	return entity
}
