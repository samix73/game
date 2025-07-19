package ecs

type ComponentID uint64

type IComponent interface {
	ID() ComponentID
	Init(id ComponentID)
}

type BaseComponent struct {
	id ComponentID
}

func (c *BaseComponent) SetID(id ComponentID) {
	if c.id != 0 {
		panic("Component ID already set")
	}
	c.id = id
}

func (c *BaseComponent) ID() ComponentID {
	return c.id
}
