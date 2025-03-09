package ecs

type EntityID uint64

type Entity struct {
	id         EntityID
	world      *World
	components map[ComponentTypeID]ComponentID
}

func NewEntity(world *World, componentTypes ...IComponentType) *Entity {
	entity := &Entity{
		id:         world.nextEntityID(),
		world:      world,
		components: make(map[ComponentTypeID]ComponentID),
	}

	for _, componentType := range componentTypes {
		c := componentType.New()

		entity.components[componentType.ID()] = c.ID()
	}

	return entity
}

func (e *Entity) ID() EntityID {
	return e.id
}
