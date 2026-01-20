package components

import (
	ecs "github.com/samix73/ebiten-ecs"
	"github.com/samix73/game/client/helpers"
	"golang.org/x/image/math/f64"
)

var _ ecs.Component = (*Collider)(nil)

type Collider struct {
	Bounds helpers.AABB
}

func (c *Collider) Init() {
	c.Bounds.Min[0] = 0
	c.Bounds.Min[1] = 0
	c.Bounds.Max[0] = 1
	c.Bounds.Max[1] = 1
}

func (c *Collider) Reset() {
	c.Bounds.Reset()
}

var _ ecs.Component = (*Collision)(nil)

type Collision struct {
	Entity      ecs.EntityID
	Penetration float64
	Normal      f64.Vec2
}

func (c *Collision) Reset() {}
