package systems

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/samix73/game/ecs"
	"github.com/samix73/game/game/components"
)

var _ ecs.System = (*TileSystem)(nil)

func init() {
	if err := ecs.RegisterSystem(NewTileSystem); err != nil {
		panic(err)
	}
}

type TileSystem struct {
	*ecs.BaseSystem
}

func NewTileSystem(priority int) *TileSystem {
	return &TileSystem{
		BaseSystem: ecs.NewBaseSystem(priority),
	}
}

func (t *TileSystem) validateTileMap(tm *components.TileMap) bool {
	if tm.Width <= 0 || tm.Height <= 0 || tm.TileSize <= 0 {
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
			tileImg := tm.ImageAt(x, y)
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

func (t *TileSystem) buildTileSet(em *ecs.EntityManager, entity ecs.EntityID, tm *components.TileMap) {
	renderable := ecs.MustGetComponent[components.Renderable](em, entity)

	renderable.Order = tm.Layer
	renderable.Sprite = t.buildTileSetImage(tm)
}

func (t *TileSystem) Update() error {
	em := t.EntityManager()

	for _, entity := range ecs.Query2[components.TileMap, components.Renderable](em) {
		tm := ecs.MustGetComponent[components.TileMap](em, entity)

		if !t.validateTileMap(tm) {
			continue
		}

		t.buildTileSet(em, entity, tm)
	}

	return nil
}

func (t *TileSystem) Teardown() {
}
