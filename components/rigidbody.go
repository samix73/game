package components

import "golang.org/x/image/math/f64"

// RigidBody represents a physics body with mass, velocity, and gravity.
type RigidBody struct {
	Mass     float64
	Velocity f64.Vec2
	Gravity  bool
}

// ApplyImpulse applies an impulse to the rigid body.
func (r *RigidBody) ApplyImpulse(impulse f64.Vec2) {
	if r.Mass <= 0 {
		return
	}

	r.Velocity[0] += impulse[0] / r.Mass
	r.Velocity[1] += impulse[1] / r.Mass
}

func (r *RigidBody) Init() {
	r.Mass = 1.0
	r.Velocity[0] = 0
	r.Velocity[1] = 0
	r.Gravity = true
}

func (r *RigidBody) Reset() {
	r.Mass = 1.0
	r.Velocity[0] = 0
	r.Velocity[1] = 0
	r.Gravity = false
}
