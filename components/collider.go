package components

import (
	"github.com/samix73/game/ecs"
	"github.com/samix73/game/helpers"
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
