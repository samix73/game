package components

import (
	"github.com/jakecoffman/cp"
	"github.com/samix73/game/ecs"
)

func init() {
	ecs.RegisterComponent[RigidBody]()
}

// RigidBody represents a physics body with mass, velocity, and gravity.
type RigidBody struct {
	Mass     float64
	Velocity cp.Vector
	Gravity  bool
}

// ApplyImpulse applies an impulse to the rigid body.
func (r *RigidBody) ApplyImpulse(impulse cp.Vector) {
	if r.Mass <= 0 {
		return
	}

	r.Velocity.X += impulse.X / r.Mass
	r.Velocity.Y += impulse.Y / r.Mass
}

func (rb *RigidBody) ApplyAcceleration(acceleration cp.Vector) {
	rb.Velocity.X += acceleration.X
	rb.Velocity.Y += acceleration.Y
}

func (r *RigidBody) Init() {
	r.Mass = 1.0
	r.Velocity.X = 0
	r.Velocity.Y = 0
	r.Gravity = true
}

func (r *RigidBody) Reset() {
	r.Mass = 1.0
	r.Velocity.X = 0
	r.Velocity.Y = 0
	r.Gravity = false
}
