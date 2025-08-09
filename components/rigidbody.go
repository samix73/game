package components

import "golang.org/x/image/math/f64"

// RigidBody represents a physics body with mass, velocity, and gravity.
type RigidBody struct {
	Mass     float64
	Velocity f64.Vec2
	Gravity  bool
}

func (r *RigidBody) Reset() {
	r.Mass = 1.0
	r.Velocity[0] = 0
	r.Velocity[1] = 0
	r.Gravity = false
}
