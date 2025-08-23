package systems

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/samix73/game/ecs"
	"github.com/samix73/game/game"
	"github.com/samix73/game/keys"
)

var _ ecs.System = (*PauseSystem)(nil)

type PauseSystem struct {
	*ecs.BaseSystem[*game.Game]

	tps int
}

func NewPauseSystem(priority int, entityManager *ecs.EntityManager, game *game.Game) *PauseSystem {
	return &PauseSystem{
		BaseSystem: ecs.NewBaseSystem(ecs.NextID(), priority, entityManager, game),

		tps: ebiten.TPS(),
	}
}

func (p *PauseSystem) Update() error {
	if !keys.IsPressed(keys.PauseAction) {
		return nil
	}

	if ebiten.TPS() == p.tps {
		ebiten.SetTPS(1)
	} else {
		ebiten.SetTPS(p.tps)
	}

	return nil
}

func (p *PauseSystem) Teardown() {}
