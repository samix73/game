package systems

import (
	"github.com/samix73/game/components"
	"github.com/samix73/game/ecs"
	"github.com/samix73/game/entities"
)

var _ ecs.System = (*Camera)(nil)

type Camera struct {
	*ecs.BaseSystem

	screenWidth  int
	screenHeight int
}

func NewCameraSystem(priority int, entityManager *ecs.EntityManager, screenWidth, screenHeight int) *Camera {
	return &Camera{
		BaseSystem: ecs.NewBaseSystem(ecs.NextID(), priority, entityManager),

		screenWidth:  screenWidth,
		screenHeight: screenHeight,
	}
}

func (c *Camera) createDefaultCamera() ecs.EntityID {
	return entities.NewCameraEntity(c.EntityManager(), true)
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

func (c *Camera) entitiesInView(cameraCenter *components.Transform) []ecs.EntityID {
	return nil
}

func (c *Camera) makeRenderable(entities []ecs.EntityID) {
	for _, entity := range entities {
		ecs.AddComponent[components.Renderable](c.EntityManager(), entity)
	}
}

func (c *Camera) Update() error {
	camera := c.activeCamera()
	cameraTransform := ecs.MustGetComponent[components.Transform](c.EntityManager(), camera)

	entitiesInView := c.entitiesInView(cameraTransform)

	c.makeRenderable(entitiesInView)

	return nil
}

func (c *Camera) Teardown() {}
