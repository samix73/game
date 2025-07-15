package ecs

import "fmt"

type World struct {
	componentTypes map[ComponentTypeID]*ComponentType[IComponent]
	entities       map[EntityID]*Entity
}

func NewWorld() *World {
	return &World{
		componentTypes: make(map[ComponentTypeID]*ComponentType[IComponent]),
		entities:       make(map[EntityID]*Entity),
	}
}

func (w *World) Update() error {
	for _, componentType := range w.componentTypes {
		if err := componentType.Update(); err != nil {
			return fmt.Errorf("error updating component type %d: %w", componentType.ID(), err)
		}
	}

	return nil
}

func (w *World) registerComponentType(componentType *ComponentType[IComponent]) {
	for _, ct := range w.componentTypes {
		if fmt.Sprintf("%T", ct) == fmt.Sprintf("%T", componentType) {
			panic(fmt.Errorf("ComponentType already registered: %T", componentType))
		}
	}

	componentType.SetID(ComponentTypeID(len(w.componentTypes) + 1))

	w.componentTypes[componentType.ID()] = componentType
}

func (w *World) registerEntity(entity *Entity) {
	if entity.id != 0 {
		panic("Entity ID already set")
	}
	entity.SetID(EntityID(len(w.componentTypes) + 1))
	w.entities[entity.ID()] = entity
}

func (w *World) GetComponentType(id ComponentTypeID) (*ComponentType[IComponent], bool) {
	ct, ok := w.componentTypes[id]

	return ct, ok
}
