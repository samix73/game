package systems

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/samix73/game/components"
	"github.com/samix73/game/ecs"
	"github.com/samix73/game/entities"
	"golang.org/x/image/math/f64"
)

var _ ecs.RendererSystem = (*Camera)(nil)

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
	cameraPos := cameraTransform.Vec2
	entityPos := entityTransform.Vec2

	// Calculate the camera's view bounds
	left := (cameraPos[0] - float64(c.screenWidth)/2.0) * camera.Zoom
	right := (cameraPos[0] + float64(c.screenWidth)/2.0) * camera.Zoom
	top := (cameraPos[1] - float64(c.screenHeight)/2.0) * camera.Zoom
	bottom := (cameraPos[1] + float64(c.screenHeight)/2.0) * camera.Zoom

	if entityPos[0] < left || entityPos[0] > right {
		return f64.Vec2{}, false
	}

	if entityPos[1] < top || entityPos[1] > bottom {
		return f64.Vec2{}, false
	}

	return f64.Vec2{
		(entityPos[0] - left) / camera.Zoom,
		(entityPos[1] - top) / camera.Zoom,
	}, true
}

func (c *Camera) Update() error {
	camera := c.activeCamera()
	cameraTransform := ecs.MustGetComponent[components.Transform](c.EntityManager(), camera)
	cameraComp := ecs.MustGetComponent[components.Camera](c.EntityManager(), camera)

	renderableEntities := ecs.Query2[components.Transform, components.Renderable](c.EntityManager())

	em := c.EntityManager()

	for entity := range renderableEntities {
		entityTransform := ecs.MustGetComponent[components.Transform](em, entity)

		onscreenPos, ok := c.inView(cameraComp, cameraTransform, entityTransform)
		if !ok {
			ecs.RemoveComponent[components.Render](em, entity)

			continue
		}

		render := ecs.AddComponent[components.Renderable](em, entity)
		render.GeoM.Translate(onscreenPos[0], onscreenPos[1])
	}

	return nil
}

func (c *Camera) Draw(screen *ebiten.Image) {
	em := c.EntityManager()

	for entity := range ecs.Query2[components.Render, components.Renderable](em) {
		renderable := ecs.MustGetComponent[components.Renderable](em, entity)

		screen.DrawImage(renderable.Sprite, &ebiten.DrawImageOptions{
			GeoM: renderable.GeoM,
		})
	}
}

func (c *Camera) Teardown() {}
