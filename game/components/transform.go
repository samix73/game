package components

import (
	"github.com/jakecoffman/cp"
	"github.com/samix73/game/ecs"
)

func init() {
	ecs.RegisterComponent[Transform]()
}

// Transform represents the position and rotation of an entity in 2D space.
type Transform struct {
	Position cp.Vector `hcl:"Position,optional"`
	Rotation float64   `hcl:"Rotation,optional"`
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
