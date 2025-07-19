package ecs

import (
	"fmt"
	"maps"
	"reflect"
)

type ComponentTypeID uint64

type ComponentType[T IComponent] struct {
	id              ComponentTypeID
	name            string
	world           *World
	values          map[ComponentID]T
	reflectType     reflect.Type
	nextComponentID ComponentID
}

func NewComponentType[T IComponent](world *World) *ComponentType[T] {
	var v T
	componentType := reflect.TypeOf(v)

	// Check if the concrete type is a pointer
	if componentType.Kind() != reflect.Ptr {
		panic(fmt.Sprintf("ComponentType %s must be a pointer type", componentType.Name()))
	}

	c := &ComponentType[T]{
		name:            componentType.Name(),
		world:           world,
		values:          make(map[ComponentID]T),
		reflectType:     reflect.TypeOf(v).Elem(),
		nextComponentID: 1,
	}

	var a any = c

	world.registerComponentType(a.(*ComponentType[IComponent]))

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

func (c *ComponentType[T]) getNextComponentID() ComponentID {
	id := c.nextComponentID
	c.nextComponentID++
	return id
}

func (c *ComponentType[T]) New() IComponent {
	var v T

	id := c.getNextComponentID()

	v.Init(id)
	c.values[id] = v

	return v
}

func (c *ComponentType[T]) GetComponentByID(id ComponentID) (T, bool) {
	value, exists := c.values[id]
	return value, exists
}

func (c *ComponentType[T]) GetAll() map[ComponentID]T {
	result := make(map[ComponentID]T)
	maps.Copy(result, c.values)

	return result
}
