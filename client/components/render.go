package components

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Renderable represents a 2D entity that can be rendered on the screen.
type Renderable struct {
	Sprite *ebiten.Image
	GeoM   ebiten.GeoM
	Order  int // Rendering order; lower values are rendered first
}

func (r *Renderable) Reset() {
	if r.Sprite != nil {
		r.Sprite.Deallocate()
	}
	r.GeoM.Reset()
	r.Order = 0
}

// Render marks entities to be rendered.
type Render struct{}
