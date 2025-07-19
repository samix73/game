package ecs

import (
	"fmt"
	"maps"
	"reflect"
)

type ComponentTypeID uint64

type IComponentType interface {
	ID() ComponentTypeID
	SetID(id ComponentTypeID)
	New() IComponent
	ReflectType() reflect.Type
	RemoveComponent(id ComponentID)
}

var _ IComponentType = (*ComponentType[IComponent])(nil)

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
	componentReflectType := reflect.TypeOf(v)

	// Check if the concrete type is a pointer
	if componentReflectType.Kind() != reflect.Ptr {
		panic(fmt.Sprintf("ComponentType %s must be a pointer type", componentReflectType.Name()))
	}

	componentType := &ComponentType[T]{
		name:            componentReflectType.Name(),
		world:           world,
		values:          make(map[ComponentID]T),
		reflectType:     componentReflectType.Elem(),
		nextComponentID: 1,
	}

	world.registerComponentType(componentType)

	return componentType
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

func (c *ComponentType[T]) ReflectType() reflect.Type {
	return c.reflectType
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

func (c *ComponentType[T]) RemoveComponent(id ComponentID) {
	if _, exists := c.values[id]; !exists {
		return
	}

	delete(c.values, id)
}
