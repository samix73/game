package systems

import (
	"log/slog"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	ecs "github.com/samix73/ebiten-ecs"
	"github.com/samix73/game/components"
	"github.com/samix73/game/entities"
	"golang.org/x/image/math/f64"
)

var _ ecs.DrawableSystem = (*Camera)(nil)

type Camera struct {
	*ecs.BaseSystem

	activeCamera ecs.EntityID
}

func NewCameraSystem(priority int) *Camera {
	return &Camera{
		BaseSystem: ecs.NewBaseSystem(ecs.NextID(), priority),
	}
}

func (c *Camera) createDefaultCamera() ecs.EntityID {
	cfg := c.Game().Config()
	return entities.NewCameraEntity(
		c.EntityManager(),
		true,
		float64(cfg.ScreenWidth), float64(cfg.ScreenHeight),
	)
}

func (c *Camera) getActiveCamera() ecs.EntityID {
	if c.activeCamera != ecs.UndefinedID {
		return c.activeCamera
	}

	em := c.EntityManager()

	activeCamera, ok := ecs.First(ecs.Query[components.ActiveCamera](em))
	if ok {
		c.activeCamera = activeCamera

		return activeCamera
	}

	camera, ok := ecs.First(ecs.Query[components.Camera](em))
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

	cfg := c.Game().Config()

	screenWidth := float64(cfg.ScreenWidth)
	screenHeight := float64(cfg.ScreenHeight)

	halfScreenWidth := screenWidth / 2
	halfScreenHeight := screenHeight / 2

	// Compute sprite top-left in screen space.
	screenX := (entX - camX) + halfScreenWidth - sw*0.5
	screenY := halfScreenHeight - (entY - camY) - sh*0.5

	// AABB vs screen rect
	if screenX+sw <= 0 || screenY+sh <= 0 || screenX >= screenWidth || screenY >= screenHeight {
		return f64.Vec2{}, false
	}

	return f64.Vec2{screenX, screenY}, true
}

func (c *Camera) Update() error {
	em := c.EntityManager()

	camera := c.getActiveCamera()
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

	renderables := make([]*components.Renderable, 0)

	for entity := range ecs.Query2[components.Render, components.Renderable](em) {
		renderable := ecs.MustGetComponent[components.Renderable](em, entity)

		if renderable.Sprite == nil {
			continue
		}

		renderables = append(renderables, renderable)
	}

	slices.SortStableFunc(renderables, func(a, b *components.Renderable) int {
		return a.Order - b.Order
	})

	for _, render := range renderables {
		screen.DrawImage(render.Sprite, &ebiten.DrawImageOptions{
			GeoM: render.GeoM,
		})
	}
}

func (c *Camera) Teardown() {}
