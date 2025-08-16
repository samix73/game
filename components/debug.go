package components

import (
	"image/color"

	"github.com/samix73/game/ecs"
	"github.com/samix73/game/helpers"
)

var _ ecs.Component = (*Collider)(nil)

type DebugBounds struct {
	Bounds helpers.AABB
	Color  color.RGBA
}

func (d *DebugBounds) Reset() {
	d.Bounds.Reset()
	d.Color.R = 0
	d.Color.G = 0
	d.Color.B = 0
	d.Color.A = 0
}
