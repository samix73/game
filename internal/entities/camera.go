package entities

import (
	"github.com/samix73/game/ecs"
	"github.com/samix73/game/internal/components"
	"golang.org/x/image/math/f64"
)

func NewCameraEntity(em *ecs.EntityManager, width, height int) ecs.EntityID {
	entity := em.NewEntity()
	em.AddComponent(entity, &components.Transform{
		Vec: f64.Vec2{0, 0},
		Rot: 0,
	})
	em.AddComponent(entity, &components.Camera{
		Width:  width,
		Height: height,
		Zoom:   1.0,
	})

	return entity
}
