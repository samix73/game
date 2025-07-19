package ecs

type ISystem interface {
	Update(world *World) error
	Priority() int
}

// BaseSystem provides common functionality for systems
type BaseSystem struct {
	priority int
}

func NewBaseSystem(priority int) BaseSystem {
	return BaseSystem{priority: priority}
}

func (s BaseSystem) Priority() int {
	return s.priority
}
