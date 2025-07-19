package ecs

import (
	"fmt"
	"sort"
)

type World struct {
	componentTypes map[ComponentTypeID]*ComponentType[IComponent]
	entities       map[EntityID]*Entity
	systems        []ISystem
}

func NewWorld() *World {
	return &World{
		componentTypes: make(map[ComponentTypeID]*ComponentType[IComponent]),
		entities:       make(map[EntityID]*Entity),
		systems:        make([]ISystem, 0),
	}
}

// AddSystem adds a system to the world and sorts by priority
func (w *World) AddSystem(system ISystem) {
	w.systems = append(w.systems, system)

	// Sort systems by priority (lower numbers run first)
	sort.Slice(w.systems, func(i, j int) bool {
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
		if err := system.Update(w); err != nil {
			return fmt.Errorf("error updating system: %w", err)
		}
	}

	return nil
}

// GetEntities returns all entities in the world
func (w *World) GetEntities() map[EntityID]*Entity {
	return w.entities
}

func (w *World) registerComponentType(inputComponentType *ComponentType[IComponent]) {
	for _, ct := range w.componentTypes {
		if inputComponentType.reflectType == ct.reflectType {
			panic(fmt.Errorf("ComponentType already registered: %T", inputComponentType))
		}
	}

	inputComponentType.SetID(ComponentTypeID(len(w.componentTypes) + 1))
	w.componentTypes[inputComponentType.ID()] = inputComponentType
}

func (w *World) registerEntity(entity *Entity) {
	if entity.id != 0 {
		panic("Entity ID already set")
	}
	entity.SetID(EntityID(len(w.entities) + 1))
	w.entities[entity.ID()] = entity
}

func (w *World) GetComponentType(id ComponentTypeID) (*ComponentType[IComponent], bool) {
	ct, ok := w.componentTypes[id]

	return ct, ok
}
