package systems

import (
	"context"
	"fmt"
	"runtime/trace"

	"github.com/samix73/game/components"
	"github.com/samix73/game/ecs"
	"github.com/samix73/game/helpers"
)

var _ ecs.System = (*Collision)(nil)

type collisionCandidate struct {
	id     ecs.EntityID
	bounds helpers.AABB
}

type Collision struct {
	*ecs.BaseSystem
}

func NewCollisionSystem(ctx context.Context, priority int, entityManager *ecs.EntityManager) *Collision {
	return &Collision{
		BaseSystem: ecs.NewBaseSystem(ctx, ecs.NextID(ctx), priority, entityManager),
	}
}

func (c *Collision) Teardown() {}

func (c *Collision) checkCollision(ctx context.Context, a, b collisionCandidate) bool {
	region := trace.StartRegion(ctx, "systems.Collision.CollisionCheck")
	defer region.End()

	return a.bounds.Overlaps(b.bounds)
}

func (c *Collision) Update(ctx context.Context) error {
	ctx, task := trace.NewTask(ctx, "systems.Collision.Update")
	defer task.End()

	em := c.EntityManager()

	active := make([]collisionCandidate, 0, 16)
	static := make([]collisionCandidate, 0, 1024)

	for entity := range ecs.Query2[components.Collider, components.Transform](ctx, em) {
		transform := ecs.MustGetComponent[components.Transform](ctx, em, entity)
		col := ecs.MustGetComponent[components.Collider](ctx, em, entity)

		adjustedBounds := col.Bounds.Add(transform.Position)

		if ecs.HasComponent[components.RigidBody](ctx, em, entity) {
			active = append(active, collisionCandidate{
				id:     entity,
				bounds: adjustedBounds,
			})
		} else {
			static = append(static, collisionCandidate{
				id:     entity,
				bounds: adjustedBounds,
			})
		}
	}

	if len(active) == 0 {
		return nil
	}

	// Active vs Static
	for _, a := range active {
		for _, b := range static {
			if c.checkCollision(ctx, a, b) {
				fmt.Println("active", a.id, "collides with", "static", b.id)
			}
		}
	}

	if len(active) < 2 {
		return nil
	}

	// Active vs Active
	for i := 0; i < len(active); i++ {
		for j := i + 1; j < len(active); j++ {
			if c.checkCollision(ctx, active[i], active[j]) {
				fmt.Println("active", active[i].id, "collides with", "active", active[j].id)
			}
		}
	}

	return nil
}
