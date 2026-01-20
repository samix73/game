package components

import (
	"golang.org/x/image/math/f64"
)

// Transform represents the position and rotation of an entity in 2D space.
type Transform struct {
	Position f64.Vec2
	Rotation float64
}

func (t *Transform) SetPosition(x, y float64) {
	t.Position[0] = x
	t.Position[1] = y
}

func (t *Transform) Translate(x, y float64) {
	t.Position[0] += x
	t.Position[1] += y
}

func (t *Transform) Reset() {
	t.Position[0] = 0
	t.Position[1] = 0
	t.Rotation = 0
}
