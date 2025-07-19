package components

import (
	"github.com/samix73/game/ecs"
	"golang.org/x/image/math/f64"
)

var _ ecs.IComponent = (*Transform)(nil)

type Transform struct {
	id ecs.ComponentID

	vec f64.Vec2 // Position in the world
	rot float64  // Rotation in radians
}

func (t *Transform) ID() ecs.ComponentID {
	return t.id
}

func (t *Transform) Init(id ecs.ComponentID) {
	if t.id != 0 {
		panic("Transform ID already set")
	}
	t.vec = f64.Vec2{0, 0}
	t.rot = 0.0
	t.id = id
}

// Update implements ecs.IComponent.
func (t *Transform) Update() error {
	return nil
}

func (t *Transform) Vec() f64.Vec2 {
	return t.vec
}

func (t *Transform) SetVec(vec f64.Vec2) {
	t.vec = vec
}

func (t *Transform) Rot() float64 {
	return t.rot
}

func (t *Transform) SetRot(rot float64) {
	t.rot = rot
}
