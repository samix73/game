package systems

import (
	"github.com/samix73/game/components"
	"github.com/samix73/game/ecs"
	"github.com/samix73/game/entities"
)

var _ ecs.System = (*Camera)(nil)

type Camera struct {
	*ecs.BaseSystem
}

func NewCameraSystem(priority int, entityManager *ecs.EntityManager) *Camera {
	return &Camera{
		BaseSystem: ecs.NewBaseSystem(ecs.NextID(), priority, entityManager),
	}
}

func (c *Camera) createDefaultCamera() ecs.EntityID {
	return entities.NewCameraEntity(c.EntityManager(), true, 800, 600)
}

func (c *Camera) activeCamera() ecs.EntityID {
	em := c.EntityManager()

	for camera := range ecs.Query2[components.Camera, components.ActiveCamera](em) {
		return camera
	}

	activeCamera := ecs.UndefinedID
	for camera := range ecs.Query[components.Camera](em) {
		if activeCamera == ecs.UndefinedID {
			ecs.AddComponent[components.ActiveCamera](em, camera)
			activeCamera = camera
		} else {
			ecs.RemoveComponent[components.ActiveCamera](em, camera)
		}
	}

	if activeCamera != ecs.UndefinedID {
		return activeCamera
	}

	return c.createDefaultCamera()
}

func (c *Camera) Update() error {
	_ = c.activeCamera()

	return nil
}

func (c *Camera) Teardown() {}
