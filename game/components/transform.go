package components

import (
	"github.com/jakecoffman/cp"
	"github.com/samix73/game/ecs"
)

func init() {
	if err := ecs.RegisterComponent[Transform](); err != nil {
		panic(err)
	}
}

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
