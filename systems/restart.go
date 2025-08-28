package systems

import (
	"fmt"

	ecs "github.com/samix73/ebiten-ecs"
	"github.com/samix73/game/keys"
)

var _ ecs.System = (*RestartSystem)(nil)

type RestartSystem struct {
	*ecs.BaseSystem
}

func NewRestartSystem(priority int, em *ecs.EntityManager, g *ecs.Game) *RestartSystem {
	return &RestartSystem{
		BaseSystem: ecs.NewBaseSystem(ecs.NextID(), priority, em, g),
	}
}

func (r *RestartSystem) Teardown() {}

func (r *RestartSystem) Update() error {
	if keys.IsPressed(keys.RestartAction) {
		if err := r.Game().RestartActiveWorld(); err != nil {
			return fmt.Errorf("systems.RestartSystem.Update RestartActiveWorld error: %w", err)
		}
	}

	return nil
}
