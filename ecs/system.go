package ecs

import (
	"context"
	"fmt"
	"runtime/trace"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
)

type SystemID = ID

type System interface {
	ID() SystemID
	Priority() int
	Update(ctx context.Context) error
	Teardown()
	baseSystem()
}

type RendererSystem interface {
	System
	Draw(ctx context.Context, screen *ebiten.Image)
}

type BaseSystem struct {
	id            SystemID
	priority      int
	entityManager *EntityManager
}

func NewBaseSystem(ctx context.Context, id SystemID, priority int, entityManager *EntityManager) *BaseSystem {
	region := trace.StartRegion(ctx, "ecs.NewBaseSystem")
	defer region.End()

	return &BaseSystem{
		id:            id,
		priority:      priority,
		entityManager: entityManager,
	}
}

func (s *BaseSystem) ID() SystemID {
	return s.id
}

func (s *BaseSystem) Priority() int {
	return s.priority
}

func (s *BaseSystem) EntityManager() *EntityManager {
	return s.entityManager
}

func (s *BaseSystem) baseSystem() {}

type SystemManager struct {
	systems       []System
	entityManager *EntityManager
}

func NewSystemManager(ctx context.Context, entityManager *EntityManager) *SystemManager {
	region := trace.StartRegion(ctx, "ecs.NewSystemManager")
	defer region.End()

	return &SystemManager{
		systems:       make([]System, 0),
		entityManager: entityManager,
	}
}

func (sm *SystemManager) sortSystems(ctx context.Context) {
	region := trace.StartRegion(ctx, "ecs.SystemManager.sortSystems")
	defer region.End()

	slices.SortStableFunc(sm.systems, func(a, b System) int {
		if a.Priority() < b.Priority() {
			return -1
		}

		if a.Priority() > b.Priority() {
			return 1
		}

		return 0
	})
}

func (sm *SystemManager) Add(ctx context.Context, systems ...System) {
	ctx, task := trace.NewTask(ctx, "ecs.SystemManager.Add")
	defer task.End()

	if len(systems) == 0 {
		return
	}

	sm.systems = append(sm.systems, systems...)

	sm.sortSystems(ctx)
}

func (sm *SystemManager) Remove(ctx context.Context, systemID SystemID) {
	region := trace.StartRegion(ctx, "ecs.SystemManager.Remove")
	defer region.End()

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

func (sm *SystemManager) Update(ctx context.Context) error {
	for _, system := range sm.systems {
		if err := system.Update(ctx); err != nil {
			return fmt.Errorf("error updating system %d: %w", system.ID(), err)
		}
	}

	return nil
}

func (sm *SystemManager) Draw(ctx context.Context, screen *ebiten.Image) {
	for _, system := range sm.systems {
		if system, ok := system.(RendererSystem); ok {
			system.Draw(ctx, screen)
		}
	}
}

func (sm *SystemManager) Teardown() {
	for _, system := range sm.systems {
		system.Teardown()
	}

	sm.systems = nil
}
