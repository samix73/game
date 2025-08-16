package systems

import (
	"github.com/samix73/game/ecs"
)

var _ ecs.System = (*LevelGen)(nil)

type LevelGen struct {
	*ecs.BaseSystem
}

func (l *LevelGen) Update() error {
	return nil
}

func (l *LevelGen) Teardown() {}
