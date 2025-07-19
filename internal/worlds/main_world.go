package worlds

import (
	"github.com/samix73/game/ecs"
	"github.com/samix73/game/internal/components"
)

func NewMainWorld() *ecs.World {
	w := ecs.NewWorld()

	// Register component types
	ecs.NewComponentType[*components.Transform](w)

	return w
}
