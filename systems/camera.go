package systems

import (
	"context"
	"fmt"
	"log/slog"
	"runtime/trace"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/samix73/game/components"
	"github.com/samix73/game/ecs"
	"github.com/samix73/game/entities"
	"golang.org/x/image/math/f64"
)

var _ ecs.RendererSystem = (*Camera)(nil)

type Camera struct {
	*ecs.BaseSystem

	screenWidth      float64
	halfScreenWidth  float64
	screenHeight     float64
	halfScreenHeight float64
}

func NewCameraSystem(ctx context.Context, priority int, entityManager *ecs.EntityManager, screenWidth, screenHeight int) *Camera {
	ctx, task := trace.NewTask(ctx, "systems.NewCameraSystem")
	defer task.End()

	return &Camera{
		BaseSystem: ecs.NewBaseSystem(ctx, ecs.NextID(ctx), priority, entityManager),

		screenWidth:      float64(screenWidth),
		halfScreenWidth:  float64(screenWidth) * 0.5,
		screenHeight:     float64(screenHeight),
		halfScreenHeight: float64(screenHeight) * 0.5,
	}
}

func (c *Camera) createDefaultCamera(ctx context.Context) ecs.EntityID {
	_, task := trace.NewTask(ctx, "systems.Camera.createDefaultCamera")
	defer task.End()

	return entities.NewCameraEntity(ctx, c.EntityManager(), true)
}

func (c *Camera) activeCamera(ctx context.Context, em *ecs.EntityManager) ecs.EntityID {
	ctx, task := trace.NewTask(ctx, "systems.Camera.activeCamera")
	defer task.End()

	for camera := range ecs.Query2[components.Camera, components.ActiveCamera](ctx, em) {
		return camera
	}

	activeCamera := ecs.UndefinedID
	for camera := range ecs.Query[components.Camera](ctx, em) {
		if activeCamera == ecs.UndefinedID {
			ecs.AddComponent[components.ActiveCamera](ctx, em, camera)
			activeCamera = camera
		} else {
			ecs.RemoveComponent[components.ActiveCamera](ctx, em, camera)
		}
	}

	if activeCamera != ecs.UndefinedID {
		return activeCamera
	}

	return c.createDefaultCamera(ctx)
}

// inView checks if the entity is within the camera's view and returns its on-screen position if it is.
func (c *Camera) inView(ctx context.Context, cameraTransform *components.Transform, entityTransform *components.Transform, sprite *ebiten.Image) (f64.Vec2, bool) {
	region := trace.StartRegion(ctx, "systems.Camera.inView")
	defer region.End()

	if sprite == nil {
		return f64.Vec2{}, false
	}

	camX, camY := cameraTransform.Vec2[0], cameraTransform.Vec2[1]
	entX, entY := entityTransform.Vec2[0], entityTransform.Vec2[1]

	sw := float64(sprite.Bounds().Dx())
	sh := float64(sprite.Bounds().Dy())

	// Compute sprite top-left in screen space.
	screenX := (entX - camX) + c.halfScreenWidth - sw*0.5
	screenY := (entY - camY) + c.halfScreenHeight - sh*0.5

	// AABB vs screen rect
	if screenX+sw <= 0 || screenY+sh <= 0 || screenX >= c.screenWidth || screenY >= c.screenHeight {
		return f64.Vec2{}, false
	}

	return f64.Vec2{screenX, screenY}, true
}

func (c *Camera) Update(ctx context.Context) error {
	ctx, task := trace.NewTask(ctx, "systems.Camera.Update")
	defer task.End()

	em := c.EntityManager()

	camera := c.activeCamera(ctx, em)
	cameraTransform := ecs.MustGetComponent[components.Transform](ctx, em, camera)

	for entity := range ecs.Query2[components.Transform, components.Renderable](ctx, em) {
		entityTransform := ecs.MustGetComponent[components.Transform](ctx, em, entity)
		render := ecs.MustGetComponent[components.Renderable](ctx, em, entity)
		if render.Sprite == nil {
			continue
		}

		onScreenPos, ok := c.inView(ctx, cameraTransform, entityTransform, render.Sprite)
		slog.Debug("Camera.Update",
			slog.Bool("in_view", ok),
			slog.Uint64("entity", uint64(entity)),
			slog.String("position", fmt.Sprintf("(%.2f, %.2f)",
				entityTransform.Vec2[0], entityTransform.Vec2[1])),
			slog.String("on_screen_position", fmt.Sprintf("(%.2f, %.2f)",
				onScreenPos[0], onScreenPos[1])),
			slog.String("camera_position", fmt.Sprintf("(%.2f, %.2f)",
				cameraTransform.Vec2[0], cameraTransform.Vec2[1])),
		)
		if !ok {
			ecs.RemoveComponent[components.Render](ctx, em, entity)

			continue
		}

		ecs.AddComponent[components.Render](ctx, em, entity)

		render.GeoM.SetElement(0, 2, onScreenPos[0])
		render.GeoM.SetElement(1, 2, onScreenPos[1])
	}

	return nil
}

func (c *Camera) Draw(ctx context.Context, screen *ebiten.Image) {
	ctx, task := trace.NewTask(ctx, "systems.Camera.Draw")
	defer task.End()

	em := c.EntityManager()

	for entity := range ecs.Query2[components.Render, components.Renderable](ctx, em) {
		renderable := ecs.MustGetComponent[components.Renderable](ctx, em, entity)

		if renderable.Sprite == nil {
			continue
		}

		screen.DrawImage(renderable.Sprite, &ebiten.DrawImageOptions{
			GeoM: renderable.GeoM,
		})
	}
}

func (c *Camera) Teardown() {}
