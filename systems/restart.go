package systems

import (
	"github.com/samix73/game/ecs"
	"github.com/samix73/game/game"
	"github.com/samix73/game/keys"
)

var _ ecs.System = (*RestartSystem)(nil)

type RestartSystem struct {
	*ecs.BaseSystem[*game.Game]
}

func NewRestartSystem(priority int, entityManager *ecs.EntityManager, game *game.Game) *RestartSystem {
	return &RestartSystem{
		BaseSystem: ecs.NewBaseSystem(ecs.NextID(), priority, entityManager, game),
	}
}

func (r *RestartSystem) Update() error {
	if keys.IsPressed(keys.RestartAction) {
		r.Game().Restart()
	}

	return nil
}

func (r *RestartSystem) Teardown() {}
