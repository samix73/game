package ecs

type ComponentID uint64

type IComponent interface {
	GetID() ComponentID
	GetTypeID() ComponentTypeID
}

type Component[T any] struct {
	id     ComponentID
	typeID ComponentTypeID

	value T
}

func (c *Component[T]) GetID() ComponentID {
	return c.id
}

func (c *Component[T]) GetTypeID() ComponentTypeID {
	return c.typeID
}
