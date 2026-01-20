package entities

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/samix73/game/components"
	"github.com/samix73/game/ecs"
)

const (
	tileSize = 16
)

func NewTileMapEntity(em *ecs.EntityManager, img *ebiten.Image, layer, width, height int, tiles []int) ecs.EntityID {
	entity := em.NewEntity()

	ecs.AddComponent[components.Transform](em, entity)
	ecs.AddComponent[components.Renderable](em, entity)
	tileMap := ecs.AddComponent[components.TileMap](em, entity)

	tileMap.TileSize = tileSize
	tileMap.Layer = layer
	tileMap.Width = width
	tileMap.Height = height
	tileMap.Atlas = img
	tileMap.Init()

	tileMap.Tiles = tiles

	return entity
}
