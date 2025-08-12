package systems

import (
	"context"

	"github.com/samix73/game/ecs"
)

var _ ecs.System = (*LevelGen)(nil)

type LevelGen struct {
	*ecs.BaseSystem
}

func (l *LevelGen) Update(ctx context.Context) error {
	return nil
}

func (l *LevelGen) Teardown() {}
