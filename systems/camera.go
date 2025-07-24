package systems

import (
	"github.com/samix73/game/components"
	"github.com/samix73/game/ecs"
	"github.com/samix73/game/entities"
	"golang.org/x/image/math/f64"
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

// inView checks if the entity is within the camera's view and returns its on-screen position if it is.
func (c *Camera) inView(camera *components.Camera, cameraTransform *components.Transform, entityTransform *components.Transform) (f64.Vec2, bool) {
	// TODO - Implement camera zoom and rotation handling
	return f64.Vec2{}, true
}

func (c *Camera) Update() error {
	camera := c.activeCamera()
	cameraTransform := ecs.MustGetComponent[components.Transform](c.EntityManager(), camera)
	cameraComp := ecs.MustGetComponent[components.Camera](c.EntityManager(), camera)

	renderableEntities := ecs.Query2[components.Transform, components.Renderable](c.EntityManager())

	em := c.EntityManager()

	for entity := range renderableEntities {
		entityTransform := ecs.MustGetComponent[components.Transform](em, entity)

		if onscreenPos, ok := c.inView(cameraComp, cameraTransform, entityTransform); ok {
			render := ecs.AddComponent[components.Render](em, entity)
			render.OnScreenPosition[0] = onscreenPos[0]
			render.OnScreenPosition[1] = onscreenPos[1]
		} else {
			ecs.RemoveComponent[components.Render](em, entity)
		}
	}

	return nil
}

func (c *Camera) Teardown() {}
