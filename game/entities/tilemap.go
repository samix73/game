package entities

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/samix73/game/ecs"
	"github.com/samix73/game/game/components"
)

const (
	tileSize = 16
)

func NewTileMapEntity(em *ecs.EntityManager, img *ebiten.Image, layer, width, height int, tiles []int) (ecs.EntityID, error) {
	entityID, err := em.NewEntity()
	if err != nil {
		return 0, fmt.Errorf("error creating entity: %w", err)
	}

	ecs.AddComponent[components.Transform](em, entityID)
	ecs.AddComponent[components.Renderable](em, entityID)
	tileMap, err := ecs.AddComponent[components.TileMap](em, entityID)
	if err != nil {
		return entityID, fmt.Errorf("error adding tilemap: %w", err)
	}

	tileMap.TileSize = tileSize
	tileMap.Layer = layer
	tileMap.Width = width
	tileMap.Height = height
	tileMap.Atlas = img
	tileMap.Init()

	tileMap.Tiles = tiles

	return entityID, nil
}
