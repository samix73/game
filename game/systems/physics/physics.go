package physics

import (
	"log/slog"

	"github.com/jakecoffman/cp"
	"github.com/samix73/game/ecs"
	"github.com/samix73/game/game/components"
)

var _ ecs.System = (*PhysicsSystem)(nil)

func init() {
	if err := ecs.RegisterSystem(NewPhysicsSystem); err != nil {
		panic(err)
	}
	if err := ecs.RegisterSystem(NewCollisionResolverSystem); err != nil {
		panic(err)
	}
}

type PhysicsSystem struct {
	*ecs.BaseSystem
}

func NewPhysicsSystem(priority int) *PhysicsSystem {
	return &PhysicsSystem{
		BaseSystem: ecs.NewBaseSystem(priority),
	}
}

func (p *PhysicsSystem) Update() error {
	em := p.EntityManager()

	// First, apply physics movement
	for entity := range ecs.Query2[components.RigidBody, components.Transform](em) {
		rigidBody := ecs.MustGetComponent[components.RigidBody](em, entity)
		transform := ecs.MustGetComponent[components.Transform](em, entity)

		game := p.Game()

		transform.Translate(
			rigidBody.Velocity.X*game.DeltaTime(),
			rigidBody.Velocity.Y*game.DeltaTime(),
		)

		slog.Debug("Physics.Update",
			slog.Uint64("entity", uint64(entity)),
			slog.Any("position", transform.Position),
			slog.Any("velocity", rigidBody.Velocity),
		)
	}

	return nil
}

func (p *PhysicsSystem) Teardown() {
}

// CollisionResolverSystem handles collision response
type CollisionResolverSystem struct {
	*ecs.BaseSystem
}

func NewCollisionResolverSystem(priority int) *CollisionResolverSystem {
	return &CollisionResolverSystem{
		BaseSystem: ecs.NewBaseSystem(priority),
	}
}

func (cr *CollisionResolverSystem) Update() error {
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
			cr.resolveStaticCollision(transform2, rigidbody2, cp.Vector{X: -collision.Normal.X, Y: -collision.Normal.Y}, collision.Penetration)
		}
		// If neither has rigidbody, no physics response needed
	}

	return nil
}

func (cr *CollisionResolverSystem) Teardown() {
}

func (cr *CollisionResolverSystem) resolveElasticCollision(transform1 *components.Transform, rb1 *components.RigidBody,
	transform2 *components.Transform, rb2 *components.RigidBody, normal cp.Vector, penetration float64) {

	// Separate objects first
	totalMass := rb1.Mass + rb2.Mass
	if totalMass > 0 {
		separation1 := penetration * (rb2.Mass / totalMass)
		separation2 := penetration * (rb1.Mass / totalMass)

		transform1.Translate(-normal.X*separation1, -normal.Y*separation1)
		transform2.Translate(normal.X*separation2, normal.Y*separation2)
	}

	// Calculate relative velocity
	relativeVelocity := cp.Vector{
		X: rb1.Velocity.X - rb2.Velocity.X,
		Y: rb1.Velocity.Y - rb2.Velocity.Y,
	}

	// Calculate relative velocity along the normal
	velocityAlongNormal := relativeVelocity.X*normal.X + relativeVelocity.Y*normal.Y

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
	impulse := cp.Vector{X: impulseScalar * normal.X, Y: impulseScalar * normal.Y}

	if rb1.Mass > 0 {
		rb1.Velocity.X += impulse.X / rb1.Mass
		rb1.Velocity.Y += impulse.Y / rb1.Mass
	}

	if rb2.Mass > 0 {
		rb2.Velocity.X -= impulse.X / rb2.Mass
		rb2.Velocity.Y -= impulse.Y / rb2.Mass
	}
}

func (cr *CollisionResolverSystem) resolveStaticCollision(transform *components.Transform, rb *components.RigidBody,
	normal cp.Vector, penetration float64) {

	// Separate the rigidbody from the static object
	transform.Translate(-normal.X*penetration, -normal.Y*penetration)

	// Calculate velocity along the normal
	velocityAlongNormal := rb.Velocity.X*normal.X + rb.Velocity.Y*normal.Y

	// Don't resolve if velocity is separating
	if velocityAlongNormal > 0 {
		return
	}

	// Remove velocity component along the normal (stop at collision)
	const restitution = 0.3 // Lower restitution for static collisions

	rb.Velocity.X -= (1 + restitution) * velocityAlongNormal * normal.X
	rb.Velocity.Y -= (1 + restitution) * velocityAlongNormal * normal.Y
}
