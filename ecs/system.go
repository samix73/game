package ecs

type ISystem interface {
	Update() error
	Priority() int
}

// BaseSystem provides common functionality for systems
type BaseSystem struct {
	priority int
	world    *World
}

func NewBaseSystem(world *World, priority int) BaseSystem {
	return BaseSystem{
		priority: priority,
		world:    world,
	}
}

func (s *BaseSystem) Priority() int {
	return s.priority
}

func (s *BaseSystem) World() *World {
	return s.world
}
