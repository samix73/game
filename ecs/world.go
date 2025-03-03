package ecs

type World struct {
	componentTypeID ComponentTypeID
	entityID        EntityID
	componentTypes  map[ComponentTypeID]IComponentType
	components      map[ComponentTypeID]map[ComponentID]IComponent
	entities        map[EntityID]*Entity
}

func NewWorld() *World {
	return &World{
		componentTypeID: 1,
		entityID:        1,
		componentTypes:  make(map[ComponentTypeID]IComponentType),
		components:      make(map[ComponentTypeID]map[ComponentID]IComponent),
		entities:        make(map[EntityID]*Entity),
	}
}

func (w *World) nextComponentTypeID() ComponentTypeID {
	id := w.componentTypeID
	w.componentTypeID++

	return id
}

func (w *World) nextEntityID() EntityID {
	id := w.entityID
	w.entityID++

	return id
}

func (w *World) Create(componentTypes ...IComponentType) *Entity {
	e := NewEntity(w)

	for _, componentType := range componentTypes {
		if componentType.GetID() == 0 {
			componentType.SetID(w.nextComponentTypeID())
			w.componentTypes[componentType.GetID()] = componentType
		}

		component := componentType.New()

		if _, ok := w.components[componentType.GetID()]; !ok {
			w.components[componentType.GetID()] = make(map[ComponentID]IComponent)
		}
		w.components[componentType.GetID()][component.GetID()] = component

		e.AddComponent(component)
	}

	w.entities[e.id] = e

	return e
}

func (w *World) CreateMany(count int, components ...IComponentType) {
	entities := make([]*Entity, count)
	for i := range count {
		entities[i] = w.Create(components...)
	}
}
