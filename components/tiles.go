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
}

func (t *TileMap) Reset() {
	t.Width = 0
	t.Height = 0
	t.Layer = 0
	t.Tiles = make([]int, 0)
}

func (t *TileMap) Init() {
	if t.Tiles != nil || t.Width <= 0 || t.Height <= 0 {
		return
	}

	t.Tiles = make([]int, t.Width*t.Height)
	for i := range t.Tiles {
		t.Tiles[i] = -1
	}
}

func (t *TileMap) index(x, y int) int {
	if x < 0 || x >= t.Width || y < 0 || y >= t.Height {
		return -1
	}

	return y*t.Width + x
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

var _ ecs.Component = (*TileSet)(nil)

type TileSet struct {
	Atlas    *ebiten.Image
	TileSize int
	Columns  int
	Sub      []*ebiten.Image
}

func (t *TileSet) Reset() {
	t.Atlas.Deallocate()
	t.TileSize = 0
	t.Columns = 0
	t.Sub = nil
}

func (t *TileSet) At(id int) *ebiten.Image {
	if id < 0 || id >= len(t.Sub) {
		return nil
	}

	return t.Sub[id]
}

func (t *TileSet) Init() {
	if t.Atlas == nil || t.TileSize <= 0 {
		return
	}

	w := t.Atlas.Bounds().Dx()
	h := t.Atlas.Bounds().Dy()

	t.Columns = w / t.TileSize
	rows := h / t.TileSize
	count := t.Columns * rows

	t.Sub = make([]*ebiten.Image, 0, count)
	for id := range count {
		x := (id % t.Columns) * t.TileSize
		y := (id / t.Columns) * t.TileSize

		subImage := t.Atlas.SubImage(
			image.Rect(x, y, x+t.TileSize, y+t.TileSize),
		).(*ebiten.Image)
		t.Sub = append(t.Sub, subImage)
	}
}
