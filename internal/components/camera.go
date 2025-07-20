package components

import "github.com/samix73/game/ecs"

var _ ecs.IComponent = (*Camera)(nil)

type Camera struct {
	ecs.BaseComponent

	width  int
	height int
	zoom   float64
	active bool
}

func (c *Camera) Init(id ecs.ComponentID) {
	c.BaseComponent = ecs.NewBaseComponent(id)
}

func (c *Camera) Active() bool {
	return c.active
}

func (c *Camera) SetActive(active bool) {
	c.active = active
}

func (c *Camera) Width() int {
	return c.width
}

func (c *Camera) SetWidth(width int) {
	c.width = width
}

func (c *Camera) Height() int {
	return c.height
}

func (c *Camera) SetHeight(height int) {
	c.height = height
}

func (c *Camera) Zoom() float64 {
	return c.zoom
}

func (c *Camera) SetZoom(zoom float64) {
	c.zoom = zoom
}
