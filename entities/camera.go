package entities

import (
	"github.com/samix73/game/components"
	"github.com/samix73/game/ecs"
)

func NewCameraEntity(em *ecs.EntityManager, active bool, width, height int) ecs.EntityID {
	entity := em.NewEntity()
	if active {
		ecs.AddComponent[components.ActiveCamera](em, entity)
	}
	ecs.AddComponent[components.Transform](em, entity)
	camera := ecs.AddComponent[components.Camera](em, entity)
	camera.Width = width
	camera.Height = height

	return entity
}
