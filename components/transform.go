package components

import (
	"golang.org/x/image/math/f64"
)

// Transform represents the position and rotation of an entity in 2D space.
type Transform struct {
	Vec2 f64.Vec2
	Rot  float64
}

func (t *Transform) Reset() {
	t.Vec2[0] = 0
	t.Vec2[1] = 0
	t.Rot = 0
}
