package ecs

import "reflect"

type ComponentTypeID uint64

type IComponentType interface {
	ID() ComponentTypeID
	SetID(id ComponentTypeID)
	New() IComponent
}

var _ IComponentType = (*ComponentType[IComponent])(nil)

type ComponentType[T IComponent] struct {
	id               ComponentTypeID
	name             string
	world            *World
	values           map[ComponentID]T
	nextComponentID_ ComponentID
}

func NewComponentType[T IComponent](world *World) *ComponentType[T] {
	var v T

	c := &ComponentType[T]{
		name:             reflect.TypeOf(v).Name(),
		world:            world,
		values:           make(map[ComponentID]T),
		nextComponentID_: 1,
	}

	world.RegisterComponentType(c)

	return c
}

func (c *ComponentType[T]) ID() ComponentTypeID {
	return c.id
}

func (c *ComponentType[T]) SetID(id ComponentTypeID) {
	if c.id != 0 {
		panic("ComponentType ID already set")
	}

	c.id = id
}

func (c *ComponentType[T]) New() IComponent {
	var v T

	id := ComponentID(len(c.values) + 1)
	c.values[id] = v

	return v
}
