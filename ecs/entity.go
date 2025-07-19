package ecs

import "maps"

type EntityID uint64

type Entity struct {
	id         EntityID
	world      *World
	components map[ComponentTypeID]ComponentID
}

func NewEntity(world *World, componentTypeIDs ...ComponentTypeID) *Entity {
	entity := &Entity{
		world:      world,
		components: make(map[ComponentTypeID]ComponentID),
	}

	for _, componentTypeID := range componentTypeIDs {
		entity.AddComponent(componentTypeID)
	}

	world.registerEntity(entity)

	return entity
}

func (e *Entity) ID() EntityID {
	return e.id
}

func (e *Entity) SetID(id EntityID) {
	if e.id != 0 {
		panic("Entity ID already set")
	}

	e.id = id
}

func (e *Entity) AddComponent(componentTypeID ComponentTypeID) {
	if _, exists := e.components[componentTypeID]; exists {
		panic("Component already added to entity")
	}

	componentType, ok := e.world.GetComponentType(componentTypeID)
	if !ok {
		panic("ComponentType not found in world")
	}

	e.components[componentTypeID] = componentType.New().ID()
}

// HasComponent checks if the entity has a component of the given type
func (e *Entity) HasComponent(componentTypeID ComponentTypeID) bool {
	_, exists := e.components[componentTypeID]
	return exists
}

func (e *Entity) GetComponentID(componentTypeID ComponentTypeID) (ComponentID, bool) {
	id, exists := e.components[componentTypeID]
	return id, exists
}

// RemoveComponent removes a component from the entity
func (e *Entity) RemoveComponent(componentTypeID ComponentTypeID) {
	e.world.RemoveComponent(componentTypeID, e.components[componentTypeID])
	delete(e.components, componentTypeID)
}

// GetAllComponentIDs returns all component IDs attached to this entity
func (e *Entity) GetAllComponentIDs() map[ComponentTypeID]ComponentID {
	result := make(map[ComponentTypeID]ComponentID)
	maps.Copy(result, e.components)
	return result
}

// GetEntityComponent retrieves a component of type T from this entity
func GetEntityComponent[T IComponent](e *Entity, componentTypeID ComponentTypeID) (T, bool) {
	var zero T

	if e == nil || componentTypeID == 0 {
		return zero, false
	}

	componentID, exists := e.components[componentTypeID]
	if !exists {
		return zero, false
	}

	component, ok := e.world.componentTypes[componentTypeID].GetComponentByID(componentID)
	if !ok {
		return zero, false
	}

	return component.(T), true
}
