package components

import "github.com/samix73/game/helpers"

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
