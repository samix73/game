package components

type ComponentID uint64

type Component interface {
	ID() ComponentID
}
