package ecs

import (
	"context"
	"runtime/trace"

	"github.com/hajimehoshi/ebiten/v2"
)

type World interface {
	Update(ctx context.Context) error
	Draw(ctx context.Context, screen *ebiten.Image)
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

func (w *BaseWorld) Update(ctx context.Context) error {
	ctx, task := trace.NewTask(ctx, "ecs.BaseWorld.Update")
	defer task.End()

	if err := w.SystemManager().Update(ctx); err != nil {
		return err
	}
	return nil
}

func (w *BaseWorld) Draw(ctx context.Context, screen *ebiten.Image) {
	ctx, task := trace.NewTask(ctx, "ecs.BaseWorld.Draw")
	defer task.End()

	w.SystemManager().Draw(ctx, screen)
}

func (w *BaseWorld) EntityManager() *EntityManager {
	return w.entityManager
}

func (w *BaseWorld) SystemManager() *SystemManager {
	return w.systemManager
}
