package ecs

type EntityID uint64

type Entity struct {
	world *World

	id         EntityID
	components map[ComponentTypeID]IComponent
}

func NewEntity(world *World) *Entity {
	return &Entity{
		world:      world,
		id:         world.nextEntityID(),
		components: make(map[ComponentTypeID]IComponent),
	}
}

func (e *Entity) AddComponent(component IComponent) {
	e.components[component.GetTypeID()] = component
}

func (e *Entity) GetID() EntityID {
	return e.id
}
