package components

import (
	"github.com/samix73/game/ecs"
	"github.com/samix73/game/helpers"
)

var _ ecs.Component = (*ColliderComponent)(nil)

type ColliderComponent struct {
	Bounds helpers.AABB
}

func (c *ColliderComponent) Init() {
	c.Bounds.Min[0] = 0
	c.Bounds.Min[1] = 0
	c.Bounds.Max[0] = 1
	c.Bounds.Max[1] = 1
}

func (c *ColliderComponent) Reset() {
	c.Bounds.Min[0] = 0
	c.Bounds.Min[1] = 0
	c.Bounds.Max[0] = 0
	c.Bounds.Max[1] = 0
}
