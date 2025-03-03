package ecs

import "fmt"

type IComponentType interface {
	GetName() string
	New() IComponent
	SetWorld(world *World)
	GetID() ComponentTypeID
	SetID(id ComponentTypeID)
}

type ComponentTypeID uint64

var _ IComponentType = (*ComponentType[any])(nil)

type ComponentType[T any] struct {
	world *World
	id    ComponentTypeID

	name       string
	nextID     ComponentID
	components map[ComponentID]*Component[T]
}

func NewComponentType[T any]() *ComponentType[T] {
	var v T
	c := &ComponentType[T]{
		name:       fmt.Sprintf("%T", v),
		components: make(map[ComponentID]*Component[T]),
		nextID:     1,
	}

	return c
}

func (c *ComponentType[T]) SetWorld(world *World) {
	if c.world != nil {
		panic("world already set")
	}

	c.world = world
}

func (c *ComponentType[T]) GetID() ComponentTypeID {
	return c.id
}

func (c *ComponentType[T]) SetID(id ComponentTypeID) {
	if c.id != 0 {
		panic("id already set")
	}

	c.id = id
}

func (c *ComponentType[T]) GetName() string {
	return c.name
}

func (c *ComponentType[T]) nextComponentID() ComponentID {
	id := c.nextID
	c.nextID++
	return id
}

func (c *ComponentType[T]) New() IComponent {
	var v T
	component := &Component[T]{
		id:    c.nextComponentID(),
		value: v,
	}

	c.components[component.id] = component

	return component
}
