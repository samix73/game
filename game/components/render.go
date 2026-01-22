package components

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/samix73/game/ecs"
)

func init() {
	ecs.RegisterComponent[Renderable]()
	ecs.RegisterComponent[Render]()
}

// Renderable represents a 2D entity that can be rendered on the screen.
type Renderable struct {
	Sprite *ebiten.Image `hcl:"-"`
	GeoM   ebiten.GeoM   `hcl:"-"`
	Order  int           `hcl:"Order,optional"` // Rendering order; lower values are rendered first
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
