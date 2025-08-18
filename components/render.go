package components

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Renderable represents a 2D entity that can be rendered on the screen.
type Renderable struct {
	Sprite *ebiten.Image
	GeoM   ebiten.GeoM
}

func (r *Renderable) Reset() {
	if r.Sprite != nil {
		r.Sprite.Deallocate()
	}
	r.GeoM.Reset()
}

// Render marks entities to be rendered.
type Render struct{}
