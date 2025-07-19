package ecs

type ComponentID uint64

type IComponent interface {
	ID() ComponentID
	Init(id ComponentID)
}

type BaseComponent struct {
	id ComponentID
}

func NewBaseComponent(id ComponentID) BaseComponent {
	return BaseComponent{
		id: id,
	}
}

func (c *BaseComponent) ID() ComponentID {
	return c.id
}
