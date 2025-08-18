package systems

import (
	"github.com/samix73/game/components"
	"github.com/samix73/game/ecs"
)

var _ ecs.System = (*PlayerCollision)(nil)

type PlayerCollision struct {
	*ecs.BaseSystem
}

func NewPlayerCollisionSystem(priority int, entityManager *ecs.EntityManager) *PlayerCollision {
	return &PlayerCollision{
		BaseSystem: ecs.NewBaseSystem(ecs.NextID(), priority, entityManager),
	}
}

func (c *PlayerCollision) Teardown() {}

func (c *PlayerCollision) IsPlayerColliding() bool {
	em := c.EntityManager()
	_, ok := ecs.First(ecs.Query2[components.Player, components.Collision](em))
	return ok
}

func (c *PlayerCollision) Update() error {
	c.IsPlayerColliding()
	return nil
}
