package components

import (
	ecs "github.com/samix73/ebiten-ecs"
	"github.com/samix73/game/helpers"
)

var _ ecs.Component = (*Camera)(nil)

// Camera represents the viewable area of the game world.
type Camera struct {
	Bounds helpers.AABB
	Zoom   float64
}

func (c *Camera) Init() {
	c.Zoom = 1.0
}

func (c *Camera) Reset() {
	c.Bounds.Reset()
	c.Zoom = 0
}

type ActiveCamera struct{}
