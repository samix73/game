package entities

import (
	"github.com/samix73/game/components"
	"github.com/samix73/game/ecs"
)

func NewCameraEntity(em *ecs.EntityManager, active bool) ecs.EntityID {
	entity := em.NewEntity()
	if active {
		ecs.AddComponent[components.ActiveCamera](em, entity)
	}
	ecs.AddComponent[components.Transform](em, entity)
	camera := ecs.AddComponent[components.Camera](em, entity)
	camera.Zoom = 1.0

	return entity
}
