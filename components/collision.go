package components

import (
	"image"

	"github.com/samix73/game/ecs"
)

var _ ecs.Component = (*CollisionComponent)(nil)

type CollisionComponent struct {
	Bounds image.Rectangle
}

func (c *CollisionComponent) Reset() {
	c.Bounds.Min.X = 0
	c.Bounds.Min.Y = 0
	c.Bounds.Max.X = 0
	c.Bounds.Max.Y = 0
}
