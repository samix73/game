package components

import (
	"github.com/jakecoffman/cp"
	ecs "github.com/samix73/ebiten-ecs"
)

var _ ecs.Component = (*Camera)(nil)

// Camera represents the viewable area of the game world.
type Camera struct {
	Bounds cp.BB
	Zoom   float64
}

func (c *Camera) Init() {
	c.Zoom = 1.0
}

func (c *Camera) Reset() {
	c.Bounds = cp.BB{}
	c.Zoom = 0
}

type ActiveCamera struct{}
