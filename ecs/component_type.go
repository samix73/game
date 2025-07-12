package ecs

import (
	"fmt"
	"reflect"
)

type ComponentTypeID uint64

type IComponentType interface {
	ID() ComponentTypeID
	SetID(id ComponentTypeID)
	Update() error
	New() IComponent
}

var _ IComponentType = (*ComponentType[IComponent])(nil)

type ComponentType[T IComponent] struct {
	id     ComponentTypeID
	name   string
	world  *World
	values map[ComponentID]T
}

func NewComponentType[T IComponent](world *World) *ComponentType[T] {
	var v T
	componentType := reflect.TypeOf(v)

	// Check if the concrete type is a pointer
	if componentType.Kind() != reflect.Ptr {
		panic(fmt.Sprintf("ComponentType %s must be a pointer type", componentType.Name()))
	}

	c := &ComponentType[T]{
		name:   componentType.Name(),
		world:  world,
		values: make(map[ComponentID]T),
	}

	world.registerComponentType(c)

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

func (c *ComponentType[T]) Update() error {
	for _, value := range c.values {
		if err := value.Update(); err != nil {
			return fmt.Errorf("error updating component %d: %w", value.ID(), err)
		}
	}

	return nil
}

func (c *ComponentType[T]) New() IComponent {
	var v T

	id := ComponentID(len(c.values) + 1)
	c.values[id] = v

	return v
}
