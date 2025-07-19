package ecs

import (
	"fmt"
	"sort"
)

type World struct {
	componentTypes      map[ComponentTypeID]IComponentType
	entities            map[EntityID]*Entity
	nextEntityID        EntityID
	nextComponentTypeID ComponentTypeID
	systems             []ISystem
}

func NewWorld() *World {
	return &World{
		componentTypes:      make(map[ComponentTypeID]IComponentType),
		entities:            make(map[EntityID]*Entity),
		systems:             make([]ISystem, 0),
		nextEntityID:        1,
		nextComponentTypeID: 1,
	}
}

func (w *World) getNextEntityID() EntityID {
	id := w.nextEntityID
	w.nextEntityID++

	return id
}

func (w *World) getNextComponentTypeID() ComponentTypeID {
	id := w.nextComponentTypeID
	w.nextComponentTypeID++

	return id
}

// AddSystem adds a system to the world and sorts by priority
func (w *World) AddSystem(system ISystem) {
	w.systems = append(w.systems, system)

	// Sort systems by priority (lower numbers run first)
	sort.SliceStable(w.systems, func(i, j int) bool {
		return w.systems[i].Priority() < w.systems[j].Priority()
	})
}

// RemoveSystem removes a system from the world
func (w *World) RemoveSystem(system ISystem) {
	for i, s := range w.systems {
		if s == system {
			w.systems = append(w.systems[:i], w.systems[i+1:]...)
			break
		}
	}
}

func (w *World) Update() error {
	// Update all systems in priority order
	for _, system := range w.systems {
		if err := system.Update(); err != nil {
			return fmt.Errorf("error updating system: %w", err)
		}
	}

	return nil
}

// GetEntities returns all entities in the world
func (w *World) GetEntities() map[EntityID]*Entity {
	return w.entities
}

func (w *World) registerComponentType(inputComponentType IComponentType) {
	for _, ct := range w.componentTypes {
		if inputComponentType.ReflectType() == ct.ReflectType() {
			panic(fmt.Errorf("ComponentType already registered: %T", inputComponentType))
		}
	}

	id := w.getNextComponentTypeID()

	inputComponentType.SetID(ComponentTypeID(id))
	w.componentTypes[id] = inputComponentType
}

func (w *World) registerEntity(entity *Entity) {
	if entity.id != 0 {
		panic("Entity ID already set")
	}

	id := w.getNextEntityID()

	entity.SetID(id)
	w.entities[id] = entity
}

func (w *World) GetComponentType(id ComponentTypeID) (IComponentType, bool) {
	ct, ok := w.componentTypes[id]

	return ct, ok
}

func (w *World) RemoveEntity(id EntityID) {
	entity, exists := w.entities[id]
	if !exists {
		return
	}

	// Remove all components associated with this entity
	for componentTypeID, componentID := range entity.GetAllComponentIDs() {
		if ct, ok := w.componentTypes[componentTypeID]; ok {
			// You'll need to add RemoveComponent method to ComponentType
			ct.RemoveComponent(componentID)
		}
	}

	delete(w.entities, id)
}

func (w *World) RemoveComponent(componentTypeID ComponentTypeID, componentID ComponentID) {
	if ct, exists := w.componentTypes[componentTypeID]; exists {
		ct.RemoveComponent(componentID)
	}
}
