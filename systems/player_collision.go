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

func (c *PlayerCollision) IsPlayerCollidingObstacle() bool {
	em := c.EntityManager()
	player, ok := ecs.First(ecs.Query2[components.Player, components.Collision](em))
	if !ok {
		return false
	}

	collision := ecs.MustGetComponent[components.Collision](em, player)

	if collision.Entity == ecs.UndefinedID {
		return false
	}

	return ecs.HasComponent[components.Obstacle](em, collision.Entity)
}

func (c *PlayerCollision) Update() error {
	if c.IsPlayerCollidingObstacle() {
		// c.Game().Pause()
	}

	return nil
}
