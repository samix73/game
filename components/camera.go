package components

import "github.com/samix73/game/ecs"

var _ ecs.Component = (*Camera)(nil)

type Camera struct {
	Width  int
	Height int
	Zoom   float64
}

func (c *Camera) Init() {
	c.Width = 0
	c.Height = 0
	c.Zoom = 1.0
}

func (c *Camera) Reset() {
	c.Width = 0
	c.Height = 0
	c.Zoom = 1.0
}

var _ ecs.Component = (*ActiveCamera)(nil)

type ActiveCamera struct{}

func (*ActiveCamera) Init()  {}
func (*ActiveCamera) Reset() {}
