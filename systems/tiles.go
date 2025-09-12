package systems

import (
	"hash/fnv"
	"unsafe"

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
	if tm.Width <= 0 || tm.Height <= 0 || tm.TileSize <= 0 {
		return false
	}

	if len(tm.Tiles) != tm.Width*tm.Height {
		return false
	}

	return true
}

func (t *TileSystem) tilesChecksum(tiles []int) uint64 {
	h := fnv.New64a()

	// Convert the int slice to a byte slice. This is an efficient way
	// to get a pointer to the underlying data without making a copy.
	// We use `unsafe` for this.
	// `unsafe.Slice` takes a pointer and a length, creating a byte slice.
	// `unsafe.Pointer(&slice[0])` gives a pointer to the first element.
	// `len(slice) * 8` is the total number of bytes (8 bytes per int64 on most systems).
	bytes := unsafe.Slice((*byte)(unsafe.Pointer(&tiles[0])), len(tiles)*8)

	_, err := h.Write(bytes)
	if err != nil {
		return 0
	}

	return h.Sum64()
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

	tm.SetRenderedTilesChecksum(t.tilesChecksum(tm.Tiles))

	return img
}

func (t *TileSystem) buildTileSet(entity ecs.EntityID, tm *components.TileMap) {
	em := t.EntityManager()

	renderable := ecs.MustGetComponent[components.Renderable](em, entity)

	renderable.Order = tm.Layer
	renderable.Sprite = t.buildTileSetImage(tm)
}

func (t *TileSystem) Update() error {
	em := t.EntityManager()

	for entity := range ecs.Query2[components.TileMap, components.Renderable](em) {
		tm := ecs.MustGetComponent[components.TileMap](em, entity)

		if !t.validateTileMap(tm) {
			continue
		}

		if tm.RenderedTilesChecksum() == t.tilesChecksum(tm.Tiles) {
			continue
		}

		t.buildTileSet(entity, tm)
	}

	return nil
}

func (t *TileSystem) Teardown() {}
