package systems

import (
	"github.com/hajimehoshi/ebiten/v2"
	ecs "github.com/samix73/ebiten-ecs"
	"github.com/samix73/game/components"
)

var _ ecs.System = (*TileSystem)(nil)

type TileSystem struct {
	*ecs.BaseSystem
}

func NewTileSystem(priority int) *TileSystem {
	return &TileSystem{
		BaseSystem: ecs.NewBaseSystem(ecs.NextID(), priority),
	}
}

func (t *TileSystem) validateTileMap(tm *components.TileMap) bool {
	if tm.Width <= 0 || tm.Height <= 0 {
		return false
	}

	if tm.Atlas == nil || tm.TileSize <= 0 || tm.Columns <= 0 {
		return false
	}

	if len(tm.Tiles) != tm.Width*tm.Height {
		return false
	}

	return true
}

func (t *TileSystem) buildTileSetImage(tm *components.TileMap) *ebiten.Image {
	img := ebiten.NewImage(tm.Width*tm.TileSize, tm.Height*tm.TileSize)

	for y := 0; y < tm.Height; y++ {
		for x := 0; x < tm.Width; x++ {
			tileIndex := tm.At(x, y)
			if tileIndex < 0 {
				continue
			}

			tileImg := tm.ImageAt(tileIndex)
			if tileImg == nil {
				continue
			}

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x*tm.TileSize), float64(y*tm.TileSize))
			img.DrawImage(tileImg, op)
		}
	}

	return img
}

func (t *TileSystem) buildTileSet(entity ecs.EntityID, tm *components.TileMap) {
	if !t.validateTileMap(tm) {
		return
	}

	em := t.EntityManager()

	renderable, ok := ecs.GetComponent[components.Renderable](em, entity)
	if !ok {
		renderable = ecs.AddComponent[components.Renderable](em, entity)
	}

	renderable.Sprite = t.buildTileSetImage(tm)
	renderable.Order = tm.Layer
}

func (t *TileSystem) Update() error {
	em := t.EntityManager()

	for entity := range ecs.Query[components.TileMap](em) {
		tm := ecs.MustGetComponent[components.TileMap](em, entity)

		t.buildTileSet(entity, tm)
	}

	return nil
}

func (t *TileSystem) Teardown() {}
