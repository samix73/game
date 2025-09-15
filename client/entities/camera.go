package entities

import (
	ecs "github.com/samix73/ebiten-ecs"
	"github.com/samix73/game/client/components"
)

func NewCameraEntity(em *ecs.EntityManager, active bool, cameraWidth, cameraHeight float64) ecs.EntityID {
	entity := em.NewEntity()
	if active {
		ecs.AddComponent[components.ActiveCamera](em, entity)
	}
	ecs.AddComponent[components.Transform](em, entity)
	camera := ecs.AddComponent[components.Camera](em, entity)
	camera.Zoom = 1.0
	camera.Bounds.Min[0] = 0
	camera.Bounds.Min[1] = 0
	camera.Bounds.Max[0] = cameraWidth
	camera.Bounds.Max[1] = cameraHeight

	return entity
}
