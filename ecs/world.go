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

type BaseWorld[G any] struct {
	entityManager *EntityManager
	systemManager *SystemManager
	game          G
}

func (bw *BaseWorld[G]) baseWorld() {
	panic("BaseWorld cannot be used directly, it must be embedded in a concrete World implementation")
}

func NewBaseWorld[G any](entityManager *EntityManager, systemManager *SystemManager, game G) *BaseWorld[G] {
	return &BaseWorld[G]{
		entityManager: entityManager,
		systemManager: systemManager,
		game:          game,
	}
}

func (w *BaseWorld[G]) Update() error {
	if err := w.SystemManager().Update(); err != nil {
		return err
	}
	return nil
}

func (w *BaseWorld[G]) Draw(screen *ebiten.Image) {
	w.SystemManager().Draw(screen)
}

func (w *BaseWorld[G]) EntityManager() *EntityManager {
	return w.entityManager
}

func (w *BaseWorld[G]) SystemManager() *SystemManager {
	return w.systemManager
}

func (w *BaseWorld[G]) Game() G {
	return w.game
}

func (m *BaseWorld[G]) Teardown() {
	m.SystemManager().Teardown()
	m.EntityManager().Teardown()
}
