package entities

import (
	"github.com/samix73/game/ecs"
	"github.com/samix73/game/internal/components"
)

func NewCameraEntity(em *ecs.EntityManager, width, height int) ecs.EntityID {
	entity := em.NewEntity()
	ecs.AddComponent[*components.Transform](em, entity)
	ecs.AddComponent[*components.Camera](em, entity)

	return entity
}
