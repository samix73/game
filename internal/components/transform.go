package components

import (
	"github.com/samix73/game/ecs"
	"golang.org/x/image/math/f64"
)

var _ ecs.IComponent = (*Transform)(nil)

type Transform struct {
	ecs.BaseComponent

	vec f64.Vec2 // Position in the world
	rot float64  // Rotation in radians
}

func (t *Transform) Init(id ecs.ComponentID) {
	t.SetID(id)

	t.vec = f64.Vec2{0, 0}
	t.rot = 0.0
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
