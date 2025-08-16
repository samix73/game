package ecs

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type World interface {
	Update() error
	Draw(screen *ebiten.Image)
	Teardown()

	baseWorld() // Force embedding BaseWorld
}

type BaseWorld struct {
	entityManager *EntityManager
	systemManager *SystemManager
}

func (bw *BaseWorld) baseWorld() {
	panic("BaseWorld cannot be used directly, it must be embedded in a concrete World implementation")
}

func NewBaseWorld(entityManager *EntityManager, systemManager *SystemManager) *BaseWorld {
	return &BaseWorld{
		entityManager: entityManager,
		systemManager: systemManager,
	}
}

func (w *BaseWorld) Update() error {
	if err := w.SystemManager().Update(); err != nil {
		return err
	}
	return nil
}

func (w *BaseWorld) Draw(screen *ebiten.Image) {
	w.SystemManager().Draw(screen)
}

func (w *BaseWorld) EntityManager() *EntityManager {
	return w.entityManager
}

func (w *BaseWorld) SystemManager() *SystemManager {
	return w.systemManager
}
