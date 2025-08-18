package systems

import (
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/samix73/game/components"
	"github.com/samix73/game/ecs"
	"github.com/samix73/game/entities"
	"github.com/samix73/game/helpers"
	"golang.org/x/image/math/f64"
)

var _ ecs.RendererSystem = (*Camera)(nil)

type Camera struct {
	*ecs.BaseSystem

	screenWidth      float64
	halfScreenWidth  float64
	screenHeight     float64
	halfScreenHeight float64
	activeCamera     ecs.EntityID
}

func NewCameraSystem(priority int, entityManager *ecs.EntityManager, screenWidth, screenHeight int) *Camera {
	return &Camera{
		BaseSystem: ecs.NewBaseSystem(ecs.NextID(), priority, entityManager),

		screenWidth:      float64(screenWidth),
		halfScreenWidth:  float64(screenWidth) * 0.5,
		screenHeight:     float64(screenHeight),
		halfScreenHeight: float64(screenHeight) * 0.5,
	}
}

func (c *Camera) createDefaultCamera() ecs.EntityID {
	return entities.NewCameraEntity(c.EntityManager(), true, float64(c.screenWidth), float64(c.screenHeight))
}

func (c *Camera) getActiveCamera(em *ecs.EntityManager) ecs.EntityID {
	if c.activeCamera != ecs.UndefinedID {
		return c.activeCamera
	}

	activeCamera, ok := helpers.First(ecs.Query[components.ActiveCamera](em))
	if ok {
		c.activeCamera = activeCamera

		return activeCamera
	}

	camera, ok := helpers.First(ecs.Query[components.Camera](em))
	if ok {
		ecs.AddComponent[components.ActiveCamera](em, camera)
		activeCamera = camera
	} else {
		activeCamera = c.createDefaultCamera()
	}

	if activeCamera == ecs.UndefinedID {
		return ecs.UndefinedID
	}

	c.activeCamera = activeCamera

	return activeCamera
}

// inView checks if the entity is within the camera's view and returns its on-screen position if it is.
func (c *Camera) inView(cameraTransform *components.Transform, entityTransform *components.Transform, sprite *ebiten.Image) (f64.Vec2, bool) {
	if sprite == nil {
		return f64.Vec2{}, false
	}

	camX, camY := cameraTransform.Position[0], cameraTransform.Position[1]
	entX, entY := entityTransform.Position[0], entityTransform.Position[1]

	sw := float64(sprite.Bounds().Dx())
	sh := float64(sprite.Bounds().Dy())

	// Compute sprite top-left in screen space.
	screenX := (entX - camX) + c.halfScreenWidth - sw*0.5
	screenY := c.halfScreenHeight - (entY - camY) - sh*0.5

	// AABB vs screen rect
	if screenX+sw <= 0 || screenY+sh <= 0 || screenX >= c.screenWidth || screenY >= c.screenHeight {
		return f64.Vec2{}, false
	}

	return f64.Vec2{screenX, screenY}, true
}

func (c *Camera) Update() error {
	em := c.EntityManager()

	camera := c.getActiveCamera(em)
	cameraTransform := ecs.MustGetComponent[components.Transform](em, camera)

	for entity := range ecs.Query2[components.Transform, components.Renderable](em) {
		entityTransform := ecs.MustGetComponent[components.Transform](em, entity)
		render := ecs.MustGetComponent[components.Renderable](em, entity)
		if render.Sprite == nil {
			continue
		}

		onScreenPos, ok := c.inView(cameraTransform, entityTransform, render.Sprite)
		slog.Debug("Camera.Update",
			slog.Bool("in_view", ok),
			slog.Uint64("entity", uint64(entity)),
			slog.Any("position", entityTransform.Position),
			slog.Any("on_screen_position", onScreenPos),
			slog.Any("camera_position", cameraTransform.Position),
		)
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
