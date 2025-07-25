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

func (c *Camera) activeCamera(em *ecs.EntityManager) ecs.EntityID {
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
	left := cameraPos[0] - float64(c.screenWidth)/2.0
	right := cameraPos[0] + float64(c.screenWidth)/2.0
	top := cameraPos[1] - float64(c.screenHeight)/2.0
	bottom := cameraPos[1] + float64(c.screenHeight)/2.0

	if entityPos[0] < left || entityPos[0] > right {
		return f64.Vec2{}, false
	}

	if entityPos[1] < top || entityPos[1] > bottom {
		return f64.Vec2{}, false
	}

	return f64.Vec2{
		entityPos[0] - left,
		entityPos[1] - top,
	}, true
}

func (c *Camera) Update() error {
	em := c.EntityManager()

	camera := c.activeCamera(em)
	cameraTransform := ecs.MustGetComponent[components.Transform](em, camera)
	cameraComp := ecs.MustGetComponent[components.Camera](em, camera)

	for entity := range ecs.Query2[components.Transform, components.Renderable](em) {
		entityTransform := ecs.MustGetComponent[components.Transform](em, entity)
		render := ecs.MustGetComponent[components.Renderable](em, entity)
		if render.Sprite == nil {
			continue
		}

		onScreenPos, ok := c.inView(cameraComp, cameraTransform, entityTransform)
		if !ok {
			ecs.RemoveComponent[components.Render](em, entity)

			continue
		}

		ecs.AddComponent[components.Render](em, entity)

		render.GeoM.SetElement(0, 2, onScreenPos[0])
		render.GeoM.SetElement(1, 2, onScreenPos[1])
	}

	return nil
}

func (c *Camera) Draw(screen *ebiten.Image) {
	em := c.EntityManager()

	for entity := range ecs.Query2[components.Render, components.Renderable](em) {
		renderable := ecs.MustGetComponent[components.Renderable](em, entity)

		if renderable.Sprite == nil {
			continue
		}

		screen.DrawImage(renderable.Sprite, &ebiten.DrawImageOptions{
			GeoM: renderable.GeoM,
		})
	}
}

func (c *Camera) Teardown() {}
