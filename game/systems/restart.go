package systems

import (
	ecs "github.com/samix73/ebiten-ecs"
	"github.com/samix73/game/keys"
)

var _ ecs.System = (*RestartSystem)(nil)

func init() {
	ecs.RegisterSystem(NewRestartSystem)
}

type RestartSystem struct {
	*ecs.BaseSystem
}

func NewRestartSystem(priority int) *RestartSystem {
	return &RestartSystem{
		BaseSystem: ecs.NewBaseSystem(priority),
	}
}

func (r *RestartSystem) Update() error {
	if !keys.IsPressed(keys.RestartAction) {
		return nil
	}

	game := r.Game()

	// Reload the main world to restart the game
	world, err := game.LoadWorld("main_world")
	if err != nil {
		return err
	}

	return game.SetActiveWorld(world)
}

func (r *RestartSystem) Start() error {
	return nil
}

func (r *RestartSystem) Teardown() {}
