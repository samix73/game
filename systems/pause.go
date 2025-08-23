package systems

import (
	"github.com/samix73/game/ecs"
	"github.com/samix73/game/game"
	"github.com/samix73/game/keys"
)

type PauseSystem struct {
	*ecs.BaseSystem[*game.Game]
}

func NewPauseSystem(priority int, entityManager *ecs.EntityManager, game *game.Game) *PauseSystem {
	return &PauseSystem{
		BaseSystem: ecs.NewBaseSystem(ecs.NextID(), priority, entityManager, game),
	}
}

func (p *PauseSystem) Update(world ecs.World) error {
	if !keys.IsPressed(keys.PauseAction) {
		return nil
	}

	if p.Game().IsPaused() {
		p.Game().Pause()
	} else {
		p.Game().Resume()
	}

	return nil
}
