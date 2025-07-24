package components

import (
	"golang.org/x/image/math/f64"
)

type Transform struct {
	Vec f64.Vec2
	Rot float64
}

func (t *Transform) Reset() {
	t.Vec[0] = 0
	t.Vec[1] = 0
	t.Rot = 0
}
