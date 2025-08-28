package systems

import (
	ecs "github.com/samix73/ebiten-ecs"
	"github.com/samix73/game/keys"
)

var _ ecs.System = (*PauseSystem)(nil)

type PauseSystem struct {
	*ecs.BaseSystem

	paused            bool
	originalTimeScale float64
}

func NewPauseSystem(priority int, entityManager *ecs.EntityManager, game *ecs.Game) *PauseSystem {
	return &PauseSystem{
		BaseSystem: ecs.NewBaseSystem(ecs.NextID(), priority, entityManager, game),

		paused: false,
	}
}

func (p *PauseSystem) Update() error {
	if !keys.IsPressed(keys.PauseAction) {
		return nil
	}

	game := p.Game()

	if p.paused {
		game.SetTimeScale(p.originalTimeScale)
	} else {
		p.originalTimeScale = game.TimeScale()
		game.SetTimeScale(0)
	}

	p.paused = !p.paused

	return nil
}

func (p *PauseSystem) Teardown() {}
