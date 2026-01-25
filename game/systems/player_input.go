package systems

import (
	"github.com/jakecoffman/cp"
	"github.com/samix73/game/ecs"
	"github.com/samix73/game/game/components"
	"github.com/samix73/game/keys"
)

var _ ecs.System = (*PlayerInputSystem)(nil)

func init() {
	if err := ecs.RegisterSystem(NewPlayerInputSystem); err != nil {
		panic(err)
	}
}

type PlayerInputSystem struct {
	*ecs.BaseSystem
}

func NewPlayerInputSystem(priority int) *PlayerInputSystem {
	return &PlayerInputSystem{
		BaseSystem: ecs.NewBaseSystem(priority),
	}
}

func (p *PlayerInputSystem) Update() error {
	if !keys.IsPressed(keys.JumpAction) {
		return nil
	}

	em := p.EntityManager()

	// Find the player entity
	for entity := range ecs.Query2[components.Player, components.RigidBody](em) {
		rb := ecs.MustGetComponent[components.RigidBody](em, entity)

		// Apply upward impulse for jump
		jumpForce := cp.Vector{X: 0, Y: 400}
		rb.ApplyImpulse(jumpForce)
	}

	return nil
}

func (p *PlayerInputSystem) Teardown() {
}
