package entities

import (
	"fmt"

	"github.com/samix73/game/assets"
	"github.com/samix73/game/components"
	"github.com/samix73/game/ecs"
)

func NewBiogEntity(em *ecs.EntityManager) (ecs.EntityID, error) {
	img, err := assets.GetSprite("biog.png")
	if err != nil {
		return ecs.UndefinedID, fmt.Errorf("error getting sprite: %v", err)
	}

	entity := em.NewEntity()
	ecs.AddComponent[components.Transform](em, entity)
	ecs.AddComponent[components.Player](em, entity)
	collider := ecs.AddComponent[components.Collider](em, entity)
	collider.Bounds.SetSize(float64(img.Bounds().Dx()), float64(img.Bounds().Dy()))

	rigidBody := ecs.AddComponent[components.RigidBody](em, entity)
	rigidBody.Gravity = true

	renderable := ecs.AddComponent[components.Renderable](em, entity)
	renderable.Sprite = img

	return entity, nil
}
