package ecs

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/pelletier/go-toml/v2"
)

type SystemConfig struct {
	Name     string `toml:"name"`
	Priority int    `toml:"priority"`
}

type EntityComponentConfig struct {
	Name string       `toml:"name"`
	Body toml.Decoder `toml:"body"`
}

type EntityConfig struct {
	Name       string                  `toml:"name"`
	Components []EntityComponentConfig `toml:"component"`
}

type WorldConfig struct {
	Name    string         `toml:"name"`
	Systems []SystemConfig `toml:"system"`
	// Entities []EntityConfig `toml:"entity"`
}

type World struct {
	cfg WorldConfig

	systemManager *SystemManager
	entityManager *EntityManager
}

// Update updates the world by updating all systems managed by the SystemManager.
// If any system returns an error during its update, the process is halted and the error is returned.
func (w *World) Update() error {
	if err := w.SystemManager().Update(); err != nil {
		return err
	}

	return nil
}

// Draw draws the world by calling the Draw method of all systems that implement the DrawableSystem interface.
func (w *World) Draw(screen *ebiten.Image) {
	w.SystemManager().Draw(screen)
}

// SystemManager returns the SystemManager associated with the world.
func (w *World) SystemManager() *SystemManager {
	return w.systemManager
}

// EntityManager returns the EntityManager associated with the world.
func (w *World) EntityManager() *EntityManager {
	return w.entityManager
}

// Teardown tears down the world by calling the Teardown method of all systems that implement the Teardowner interface
// and then tearing down the EntityManager.
func (w *World) Teardown() {
	w.SystemManager().Teardown()
	w.EntityManager().Teardown()
}
