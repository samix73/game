package ecs

import "fmt"

type World struct {
	componentTypes       map[ComponentTypeID]IComponentType
	entities             map[EntityID]*Entity
	nextComponentTypeID_ ComponentTypeID
	nextEntityID_        EntityID
}

func NewWorld() *World {
	return &World{
		componentTypes:       make(map[ComponentTypeID]IComponentType),
		nextComponentTypeID_: 1,
		nextEntityID_:        1,
	}
}

func (w *World) RegisterComponentType(componentType IComponentType) {
	for _, ct := range w.componentTypes {
		if fmt.Sprintf("%T", ct) == fmt.Sprintf("%T", componentType) {
			panic(fmt.Errorf("ComponentType already registered: %T", componentType))
		}
	}

	componentType.SetID(w.nextComponentTypeID())

	w.componentTypes[componentType.ID()] = componentType
}

func (w *World) GetComponentType(id ComponentTypeID) (IComponentType, bool) {
	ct, ok := w.componentTypes[id]

	return ct, ok
}

func (w *World) Create(componentTypes ...IComponentType) *Entity {
	entity := NewEntity(w, componentTypes...)

	w.entities[entity.ID()] = entity

	return entity
}

func (w *World) CreateMany(count int, componentTypes ...IComponentType) []*Entity {
	entities := make([]*Entity, count)

	for i := range count {
		entities[i] = w.Create(componentTypes...)
	}

	return entities
}

func (w *World) Destroy(entity *Entity) {
	delete(w.entities, entity.ID())
}

func (w *World) nextComponentTypeID() ComponentTypeID {
	id := w.nextComponentTypeID_
	w.nextComponentTypeID_++

	return id
}

func (w *World) nextEntityID() EntityID {
	id := w.nextEntityID_
	w.nextEntityID_++

	return id
}
