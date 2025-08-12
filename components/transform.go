package components

import (
	"golang.org/x/image/math/f64"
)

// Transform represents the position and rotation of an entity in 2D space.
type Transform struct {
	position f64.Vec2
	rot      float64
}

func (t *Transform) Position() f64.Vec2 {
	return t.position
}

func (t *Transform) Rotation() float64 {
	return t.rot
}

func (t *Transform) SetPosition(v f64.Vec2) {
	t.position[0] = v[0]
	t.position[1] = v[1]
}

func (t *Transform) Translate(v f64.Vec2) {
	t.position[0] += v[0]
	t.position[1] += v[1]
}

func (t *Transform) Reset() {
	t.position[0] = 0
	t.position[1] = 0
	t.rot = 0
}
