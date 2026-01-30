package physics

import (
	"github.com/jakecoffman/cp"
	ecs "github.com/samix73/ebiten-ecs"
	"github.com/samix73/game/game/components"
)

var gravity = cp.Vector{X: 0, Y: -981}

var _ ecs.System = (*GravitySystem)(nil)

func init() {
	ecs.RegisterSystem(NewGravitySystem)
}

type GravitySystem struct {
	*ecs.BaseSystem

	dv cp.Vector
}

func NewGravitySystem(priority int) *GravitySystem {
	return &GravitySystem{
		BaseSystem: ecs.NewBaseSystem(priority),
		dv:         gravity,
	}
}

func (g *GravitySystem) Update() error {
	em := g.EntityManager()
	for _, entity := range ecs.Query[components.RigidBody](em) {
		rigidBody := ecs.MustGetComponent[components.RigidBody](em, entity)
		if rigidBody == nil {
			continue
		}

		if !rigidBody.Gravity {
			continue
		}

		game := g.Game()

		acc := cp.Vector{
			X: g.dv.X * game.DeltaTime(),
			Y: g.dv.Y * game.DeltaTime(),
		}

		rigidBody.ApplyAcceleration(acc)
	}

	return nil
}

func (g *GravitySystem) Start() error {
	return nil
}

func (g *GravitySystem) Teardown() {
}
