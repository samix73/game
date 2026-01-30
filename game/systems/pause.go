package systems

import (
	"log/slog"

	ecs "github.com/samix73/ebiten-ecs"
	"github.com/samix73/game/keys"
)

var _ ecs.System = (*PauseSystem)(nil)

func init() {
	ecs.RegisterSystem(NewPauseSystem)
}

type PauseSystem struct {
	*ecs.BaseSystem

	paused            bool
	originalTimeScale float64
}

func NewPauseSystem(priority int) *PauseSystem {
	return &PauseSystem{
		BaseSystem: ecs.NewBaseSystem(priority),

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

	slog.Info("Paused", "paused", p.paused)

	return nil
}

func (p *PauseSystem) Start() error {
	return nil
}

func (p *PauseSystem) Teardown() {
}
