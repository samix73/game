package systems

import (
	"github.com/samix73/game/components"
	"github.com/samix73/game/ecs"
	"github.com/samix73/game/game"
)

var _ ecs.System = (*PlayerCollision)(nil)

type PlayerCollision struct {
	*ecs.BaseSystem[*game.Game]
}

func NewPlayerCollisionSystem(priority int, entityManager *ecs.EntityManager, game *game.Game) *PlayerCollision {
	return &PlayerCollision{
		BaseSystem: ecs.NewBaseSystem(ecs.NextID(), priority, entityManager, game),
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
	c.Game()
	return nil
}
