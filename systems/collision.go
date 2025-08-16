package systems

import (
	"context"
	"runtime/trace"

	"github.com/samix73/game/components"
	"github.com/samix73/game/ecs"
)

var _ ecs.System = (*Collision)(nil)

type Collision struct {
	*ecs.BaseSystem
}

func NewCollisionSystem(ctx context.Context, priority int, entityManager *ecs.EntityManager) *Collision {
	return &Collision{
		BaseSystem: ecs.NewBaseSystem(ctx, ecs.NextID(ctx), priority, entityManager),
	}
}

func (c *Collision) Teardown() {}

func (c *Collision) checkCollision(ctx context.Context, a, b ecs.EntityID) bool {
	region := trace.StartRegion(ctx, "systems.Collision.checkCollision")
	defer region.End()

	em := c.EntityManager()

	aCollider := ecs.MustGetComponent[components.ColliderComponent](ctx, em, a)
	bCollider := ecs.MustGetComponent[components.ColliderComponent](ctx, em, b)

	return aCollider.Bounds.Overlaps(bCollider.Bounds)
}

func (c *Collision) moveCollider(ctx context.Context, entity ecs.EntityID) {
	region := trace.StartRegion(ctx, "systems.Collision.moveCollider")
	defer region.End()

	em := c.EntityManager()
	transform := ecs.MustGetComponent[components.Transform](ctx, em, entity)
	collider := ecs.MustGetComponent[components.ColliderComponent](ctx, em, entity)

	collider.Bounds.Add(transform.Position)
}

func (c *Collision) Update(ctx context.Context) error {
	ctx, task := trace.NewTask(ctx, "systems.Collision.Update")
	defer task.End()

	collisionCandidates := make([]ecs.EntityID, 0)
	for entity := range ecs.Query2[components.ColliderComponent, components.Transform](ctx, c.EntityManager()) {
		c.moveCollider(ctx, entity)

		collisionCandidates = append(collisionCandidates, entity)
	}

	if len(collisionCandidates) < 2 {
		return nil
	}

	for i := range collisionCandidates {
		for j := i + 1; j < len(collisionCandidates); j++ {
			if c.checkCollision(ctx, collisionCandidates[i], collisionCandidates[j]) {

			}
		}
	}

	return nil
}
