package components

import (
	"github.com/jakecoffman/cp"
	"github.com/samix73/game/ecs"
)

func init() {
	ecs.RegisterComponent[Camera]()
	ecs.RegisterComponent[ActiveCamera]()
}

var _ ecs.Component = (*Camera)(nil)

// Camera represents the viewable area of the game world.
type Camera struct {
	Bounds cp.BB   `hcl:"Bounds,optional"`
	Zoom   float64 `hcl:"Zoom,optional"`
}

func (c *Camera) Init() {
	c.Zoom = 1.0
}

func (c *Camera) Reset() {
	c.Bounds = cp.BB{}
	c.Zoom = 0
}

type ActiveCamera struct{}
