package entities

import (
	"github.com/jakecoffman/cp"
	"github.com/samix73/game/components"
	"github.com/samix73/game/ecs"
)

func NewCameraEntity(em *ecs.EntityManager, active bool, cameraWidth, cameraHeight float64) ecs.EntityID {
	entity := em.NewEntity()
	if active {
		ecs.AddComponent[components.ActiveCamera](em, entity)
	}
	ecs.AddComponent[components.Transform](em, entity)
	camera := ecs.AddComponent[components.Camera](em, entity)
	camera.Zoom = 1.0
	camera.Bounds = cp.BB{L: 0, B: 0, R: cameraWidth, T: cameraHeight}

	return entity
}
