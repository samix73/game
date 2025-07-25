package components

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Renderable struct {
	Sprite *ebiten.Image
	GeoM   ebiten.GeoM
}

func (r *Renderable) Reset() {
	r.Sprite.Deallocate()
	r.GeoM.Reset()
}

type Render struct{}
