package systems

import (
	"log/slog"

	ecs "github.com/samix73/ebiten-ecs"
	"github.com/samix73/game/client/components"
	"golang.org/x/image/math/f64"
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

	// First, apply physics movement
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

// CollisionResolver handles collision response
type CollisionResolver struct {
	*ecs.BaseSystem
}

func NewCollisionResolverSystem(priority int) *CollisionResolver {
	return &CollisionResolver{
		BaseSystem: ecs.NewBaseSystem(ecs.NextID(), priority),
	}
}

func (cr *CollisionResolver) Update() error {
	em := cr.EntityManager()

	// Handle all collision responses
	for entity := range ecs.Query[components.Collision](em) {
		collision := ecs.MustGetComponent[components.Collision](em, entity)
		otherEntity := collision.Entity

		// Get components for both entities
		transform1, hasTransform1 := ecs.GetComponent[components.Transform](em, entity)
		transform2, hasTransform2 := ecs.GetComponent[components.Transform](em, otherEntity)

		if !hasTransform1 || !hasTransform2 {
			continue
		}

		rigidbody1, hasRigidBody1 := ecs.GetComponent[components.RigidBody](em, entity)
		rigidbody2, hasRigidBody2 := ecs.GetComponent[components.RigidBody](em, otherEntity)

		if collision.Penetration <= 0 {
			continue // No collision
		}

		// Resolve collision based on rigidbody presence
		if hasRigidBody1 && hasRigidBody2 {
			// Both have rigidbodies - elastic collision with mass consideration
			cr.resolveElasticCollision(transform1, rigidbody1, transform2, rigidbody2, collision.Normal, collision.Penetration)
		} else if hasRigidBody1 && !hasRigidBody2 {
			// Entity 1 has rigidbody, entity 2 is static
			cr.resolveStaticCollision(transform1, rigidbody1, collision.Normal, collision.Penetration)
		} else if !hasRigidBody1 && hasRigidBody2 {
			// Entity 2 has rigidbody, entity 1 is static
			cr.resolveStaticCollision(transform2, rigidbody2, f64.Vec2{-collision.Normal[0], -collision.Normal[1]}, collision.Penetration)
		}
		// If neither has rigidbody, no physics response needed
	}

	return nil
}

func (cr *CollisionResolver) resolveElasticCollision(transform1 *components.Transform, rb1 *components.RigidBody,
	transform2 *components.Transform, rb2 *components.RigidBody, normal f64.Vec2, penetration float64) {

	// Separate objects first
	totalMass := rb1.Mass + rb2.Mass
	if totalMass > 0 {
		separation1 := penetration * (rb2.Mass / totalMass)
		separation2 := penetration * (rb1.Mass / totalMass)

		transform1.Translate(-normal[0]*separation1, -normal[1]*separation1)
		transform2.Translate(normal[0]*separation2, normal[1]*separation2)
	}

	// Calculate relative velocity
	relativeVelocity := f64.Vec2{
		rb1.Velocity[0] - rb2.Velocity[0],
		rb1.Velocity[1] - rb2.Velocity[1],
	}

	// Calculate relative velocity along the normal
	velocityAlongNormal := relativeVelocity[0]*normal[0] + relativeVelocity[1]*normal[1]

	// Don't resolve if velocities are separating
	if velocityAlongNormal > 0 {
		return
	}

	// Restitution (bounciness) - can be made configurable later
	const restitution = 0.8

	// Calculate impulse scalar
	impulseScalar := -(1 + restitution) * velocityAlongNormal
	if rb1.Mass > 0 && rb2.Mass > 0 {
		impulseScalar /= (1/rb1.Mass + 1/rb2.Mass)
	}

	// Apply impulse
	impulse := f64.Vec2{impulseScalar * normal[0], impulseScalar * normal[1]}

	if rb1.Mass > 0 {
		rb1.Velocity[0] += impulse[0] / rb1.Mass
		rb1.Velocity[1] += impulse[1] / rb1.Mass
	}

	if rb2.Mass > 0 {
		rb2.Velocity[0] -= impulse[0] / rb2.Mass
		rb2.Velocity[1] -= impulse[1] / rb2.Mass
	}
}

func (cr *CollisionResolver) resolveStaticCollision(transform *components.Transform, rb *components.RigidBody,
	normal f64.Vec2, penetration float64) {

	// Separate the rigidbody from the static object
	transform.Translate(-normal[0]*penetration, -normal[1]*penetration)

	// Calculate velocity along the normal
	velocityAlongNormal := rb.Velocity[0]*normal[0] + rb.Velocity[1]*normal[1]

	// Don't resolve if velocity is separating
	if velocityAlongNormal > 0 {
		return
	}

	// Remove velocity component along the normal (stop at collision)
	const restitution = 0.3 // Lower restitution for static collisions

	rb.Velocity[0] -= (1 + restitution) * velocityAlongNormal * normal[0]
	rb.Velocity[1] -= (1 + restitution) * velocityAlongNormal * normal[1]
}
