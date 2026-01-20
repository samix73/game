package components

import "github.com/jakecoffman/cp"

// Transform represents the position and rotation of an entity in 2D space.
type Transform struct {
	Position cp.Vector
	Rotation float64
}

func (t *Transform) SetPosition(x, y float64) {
	t.Position.X = x
	t.Position.Y = y
}

func (t *Transform) Translate(x, y float64) {
	t.Position.X += x
	t.Position.Y += y
}

func (t *Transform) Reset() {
	t.Position.X = 0
	t.Position.Y = 0
	t.Rotation = 0
}
