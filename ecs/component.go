package ecs

type ComponentID uint64

type IComponent interface {
	ID() ComponentID
}
