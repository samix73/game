package components

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Renderable struct {
	Sprite *ebiten.Image
	GeoM   ebiten.GeoM
}

func (r *Renderable) Reset() {
	r.Sprite = new(ebiten.Image)
	r.GeoM.Reset()
}

type Render struct{}
