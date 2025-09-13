package systems

import (
	"log/slog"

	ecs "github.com/samix73/ebiten-ecs"
	"github.com/samix73/game/components"
)

type Physics struct {
	*ecs.BaseSystem
}

func NewPhysicsSystem(priority int) *Physics {
	return &Physics{
		BaseSystem: ecs.NewBaseSystem(ecs.NextID(), priority),
	}
}

func (p *Physics) Update() error {
	em := p.EntityManager()

	for entity := range ecs.Query2[components.RigidBody, components.Transform](em) {
		rigidBody := ecs.MustGetComponent[components.RigidBody](em, entity)
		transform := ecs.MustGetComponent[components.Transform](em, entity)

		game := p.Game()

		transform.Translate(
			rigidBody.Velocity[0]*game.DeltaTime(),
			rigidBody.Velocity[1]*game.DeltaTime(),
		)

		slog.Debug("Physics.Update",
			slog.Uint64("entity", uint64(entity)),
			slog.Any("position", transform.Position),
			slog.Any("velocity", rigidBody.Velocity),
		)
	}

	return nil
}
