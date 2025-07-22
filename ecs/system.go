package ecs

import (
	"fmt"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
)

type SystemID = ID

type System interface {
	ID() SystemID
	Priority() int
	Update() error
	Draw(screen *ebiten.Image)
	Teardown()
}

type SystemManager struct {
	systems       []System
	entityManager *EntityManager
}

func NewSystemManager(entityManager *EntityManager) *SystemManager {
	return &SystemManager{
		systems:       make([]System, 0),
		entityManager: entityManager,
	}
}

func (sm *SystemManager) Add(system System) {
	sm.systems = append(sm.systems, system)

	slices.SortFunc(sm.systems, func(a, b System) int {
		if a.Priority() < b.Priority() {
			return -1
		}

		if a.Priority() > b.Priority() {
			return 1
		}

		return 0
	})
}

func (sm *SystemManager) Remove(systemID SystemID) {
	indexToDelete, exists := slices.BinarySearchFunc(sm.systems, systemID, func(s System, id SystemID) int {
		if s.ID() < id {
			return -1
		}

		if s.ID() > id {
			return 1
		}

		return 0
	})

	if !exists {
		return
	}

	systemToDelete := sm.systems[indexToDelete]
	sm.systems[indexToDelete] = sm.systems[len(sm.systems)-1]
	sm.systems = sm.systems[:len(sm.systems)-1]

	systemToDelete.Teardown()
}

func (sm *SystemManager) Update() error {
	for _, system := range sm.systems {
		if err := system.Update(); err != nil {
			return fmt.Errorf("error updating system %d: %w", system.ID(), err)
		}
	}

	return nil
}

func (sm *SystemManager) Draw(screen *ebiten.Image) {
	for _, system := range sm.systems {
		system.Draw(screen)
	}
}

func (sm *SystemManager) Teardown() {
	for _, system := range sm.systems {
		system.Teardown()
	}

	sm.systems = nil
}
