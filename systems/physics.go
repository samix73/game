package systems

import (
	"log/slog"

	"github.com/samix73/game/components"
	"github.com/samix73/game/ecs"
	"github.com/samix73/game/game"
)

type Physics struct {
	*ecs.BaseSystem[*game.Game]
}

func NewPhysicsSystem(priority int, entityManager *ecs.EntityManager, game *game.Game) *Physics {
	return &Physics{
		BaseSystem: ecs.NewBaseSystem(ecs.NextID(), priority, entityManager, game),
	}
}

func (p *Physics) Teardown() {}

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
