package components

import (
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	ecs "github.com/samix73/ebiten-ecs"
)

var _ ecs.Component = (*TileMap)(nil)

type TileMap struct {
	Width, Height int
	Layer         int
	Tiles         []int // Width * Height; each int is an index into the tileset; -1 = empty
	TileSize      int

	renderedTilesChecksum uint64 // hash of the Tiles slice to detect changes
	atlas                 *ebiten.Image
	columns               int
	sub                   []*ebiten.Image
}

func (t *TileMap) Reset() {
	t.Width = 0
	t.Height = 0
	t.Layer = 0
	t.Tiles = make([]int, 0)
	t.TileSize = 0
	t.renderedTilesChecksum = 0
	t.atlas.Deallocate()
	t.columns = 0
	t.sub = nil
}

func (t *TileMap) Init() {
	if t.atlas == nil || t.TileSize <= 0 {
		return
	}

	if t.Tiles != nil || t.Width <= 0 || t.Height <= 0 {
		return
	}

	t.Tiles = make([]int, t.Width*t.Height)
	for i := range t.Tiles {
		t.Tiles[i] = -1
	}

	w := t.atlas.Bounds().Dx()
	h := t.atlas.Bounds().Dy()

	t.columns = w / t.TileSize
	rows := h / t.TileSize
	count := t.columns * rows

	t.sub = make([]*ebiten.Image, 0, count)
	for id := range count {
		x := (id % t.columns) * t.TileSize
		y := (id / t.columns) * t.TileSize

		subImage := t.atlas.SubImage(
			image.Rect(x, y, x+t.TileSize, y+t.TileSize),
		).(*ebiten.Image)
		t.sub = append(t.sub, subImage)
	}

	t.columns = t.atlas.Bounds().Dx() / t.TileSize
}

func (t *TileMap) index(x, y int) int {
	if x < 0 || x >= t.Width || y < 0 || y >= t.Height {
		return -1
	}

	return y*t.Width + x
}

func (t *TileMap) RenderedTilesChecksum() uint64 {
	return t.renderedTilesChecksum
}

func (t *TileMap) SetRenderedTilesChecksum(c uint64) {
	t.renderedTilesChecksum = c
}

func (t *TileMap) SetAtlas(atlas *ebiten.Image) {
	t.atlas = atlas
	t.Init()
}

func (t *TileMap) At(x, y int) int {
	i := t.index(x, y)
	if i == -1 {
		return -1
	}

	return t.Tiles[i]
}

func (t *TileMap) Set(x, y, id int) error {
	i := t.index(x, y)
	if i == -1 {
		return fmt.Errorf("invalid tile coordinates: (%d, %d)", x, y)
	}

	t.Tiles[i] = id

	return nil
}

func (t *TileMap) ImageAt(id int) *ebiten.Image {
	if id < 0 || id >= len(t.sub) {
		return nil
	}

	return t.sub[id]
}
