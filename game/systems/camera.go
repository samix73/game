package systems

import (
	"fmt"
	"log/slog"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jakecoffman/cp"
	ecs "github.com/samix73/ebiten-ecs"
	"github.com/samix73/game/game/components"
	"github.com/samix73/game/game/entities"
	"github.com/samix73/game/helpers"
)

var _ ecs.DrawableSystem = (*CameraSystem)(nil)

func init() {
	ecs.RegisterSystem(NewCameraSystem)
}

type CameraSystem struct {
	*ecs.BaseSystem

	activeCamera ecs.EntityID
}

func NewCameraSystem(priority int) *CameraSystem {
	return &CameraSystem{
		BaseSystem: ecs.NewBaseSystem(priority),
	}
}

func (c *CameraSystem) createDefaultCamera() (ecs.EntityID, error) {
	cfg := c.Game().Config()
	entityID, err := entities.NewCameraEntity(
		c.EntityManager(),
		true,
		float64(cfg.ScreenWidth), float64(cfg.ScreenHeight),
	)
	if err != nil {
		return 0, fmt.Errorf("error creating default camera: %w", err)
	}

	return entityID, nil
}

func (c *CameraSystem) getActiveCamera() (ecs.EntityID, error) {
	if c.activeCamera != ecs.UndefinedSystemID {
		return c.activeCamera, nil
	}

	em := c.EntityManager()

	activeCamera, ok := helpers.First(ecs.Query[components.ActiveCamera](em))
	if ok {
		c.activeCamera = activeCamera

		return activeCamera, nil
	}

	camera, ok := helpers.First(ecs.Query[components.Camera](em))
	if ok {
		ecs.AddComponent[components.ActiveCamera](em, camera)
		activeCamera = camera
	} else {
		activeCamera, err := c.createDefaultCamera()
		if err != nil {
			return 0, fmt.Errorf("error creating default camera: %w", err)
		}

		c.activeCamera = activeCamera
	}

	return activeCamera, nil
}

// inView checks if the entity is within the camera's view and returns its on-screen position if it is.
func (c *CameraSystem) inView(cameraTransform *components.Transform, entityTransform *components.Transform, sprite *ebiten.Image) (cp.Vector, bool) {
	if sprite == nil {
		return cp.Vector{}, false
	}

	camX, camY := cameraTransform.Position.X, cameraTransform.Position.Y
	entX, entY := entityTransform.Position.X, entityTransform.Position.Y

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
		return cp.Vector{}, false
	}

	return cp.Vector{X: screenX, Y: screenY}, true
}

func (c *CameraSystem) Update() error {
	em := c.EntityManager()

	camera, err := c.getActiveCamera()
	if err != nil {
		return fmt.Errorf("error getting active camera: %w", err)
	}

	cameraTransform := ecs.MustGetComponent[components.Transform](em, camera)

	for _, entity := range ecs.Query2[components.Transform, components.Renderable](em) {
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

		if !ecs.HasComponent[components.Render](em, entity) {
			if _, err := ecs.AddComponent[components.Render](em, entity); err != nil {
				return fmt.Errorf("systems.CameraSystem.Update error adding render component: %w", err)
			}
		}

		render.GeoM.SetElement(0, 2, onScreenPos.X)
		render.GeoM.SetElement(1, 2, onScreenPos.Y)
	}

	return nil
}

func (c *CameraSystem) Draw(screen *ebiten.Image) {
	em := c.EntityManager()

	renderables := make([]*components.Renderable, 0)

	for _, entity := range ecs.Query2[components.Render, components.Renderable](em) {
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

func (c *CameraSystem) Start() error {
	return nil
}

func (c *CameraSystem) Teardown() {
}
