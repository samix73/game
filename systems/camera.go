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

	screenWidth  int
	screenHeight int
}

func NewCameraSystem(ctx context.Context, priority int, entityManager *ecs.EntityManager, screenWidth, screenHeight int) *Camera {
	ctx, task := trace.NewTask(ctx, "systems.NewCameraSystem")
	defer task.End()

	return &Camera{
		BaseSystem: ecs.NewBaseSystem(ctx, ecs.NextID(ctx), priority, entityManager),

		screenWidth:  screenWidth,
		screenHeight: screenHeight,
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
func (c *Camera) inView(ctx context.Context, camera *components.Camera, cameraTransform *components.Transform, entityTransform *components.Transform, sprite *ebiten.Image) (f64.Vec2, bool) {
	ctx, task := trace.NewTask(ctx, "systems.Camera.inView")
	defer task.End()

	cameraPos := cameraTransform.Vec2
	entityPos := entityTransform.Vec2

	// Calculate the camera's view bounds
	left := (cameraPos[0] - float64(c.screenWidth)/2.0) / camera.Zoom
	right := (cameraPos[0] + float64(c.screenWidth)/2.0) / camera.Zoom
	top := (cameraPos[1] - float64(c.screenHeight)/2.0) / camera.Zoom
	bottom := (cameraPos[1] + float64(c.screenHeight)/2.0) / camera.Zoom

	if entityPos[0] < left || entityPos[0] > right {
		return f64.Vec2{}, false
	}

	if entityPos[1] < top || entityPos[1] > bottom {
		return f64.Vec2{}, false
	}

	// Calculate the on-screen position of the entity
	spriteWidth := sprite.Bounds().Dx()
	spriteHeight := sprite.Bounds().Dy()

	return f64.Vec2{
		(entityPos[0] - float64(spriteWidth/2)) - left,
		(entityPos[1] - float64(spriteHeight/2)) - top,
	}, true
}

func (c *Camera) Update(ctx context.Context) error {
	ctx, task := trace.NewTask(ctx, "systems.Camera.Update")
	defer task.End()

	em := c.EntityManager()

	camera := c.activeCamera(ctx, em)
	cameraTransform := ecs.MustGetComponent[components.Transform](ctx, em, camera)
	cameraComp := ecs.MustGetComponent[components.Camera](ctx, em, camera)

	for entity := range ecs.Query2[components.Transform, components.Renderable](ctx, em) {
		entityTransform := ecs.MustGetComponent[components.Transform](ctx, em, entity)
		render := ecs.MustGetComponent[components.Renderable](ctx, em, entity)
		if render.Sprite == nil {
			continue
		}

		onScreenPos, ok := c.inView(ctx, cameraComp, cameraTransform, entityTransform, render.Sprite)
		slog.Debug("Camera.Update",
			slog.Bool("in_view", ok),
			slog.Uint64("entity", uint64(entity)),
			slog.String("position", fmt.Sprintf("(%.2f, %.2f)",
				entityTransform.Vec2[0], entityTransform.Vec2[1])),
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

		slog.Debug("Camera.Draw", slog.Uint64("entity", uint64(entity)))
		screen.DrawImage(renderable.Sprite, &ebiten.DrawImageOptions{
			GeoM: renderable.GeoM,
		})
	}
}

func (c *Camera) Teardown() {}
