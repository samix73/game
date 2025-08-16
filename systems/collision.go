package systems

import (
	"context"
	"image"
	"runtime/trace"
	"slices"

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

	aBounds := ecs.MustGetComponent[components.CollisionComponent](ctx, em, a)
	bBounds := ecs.MustGetComponent[components.CollisionComponent](ctx, em, b)
	aTransform := ecs.MustGetComponent[components.Transform](ctx, em, a)
	bTransform := ecs.MustGetComponent[components.Transform](ctx, em, b)

	adjustedABounds := aBounds.Bounds.Add(image.Point{
		X: int(aTransform.Position[0]),
		Y: int(aTransform.Position[1]),
	})
	adjustedBBounds := bBounds.Bounds.Add(image.Point{
		X: int(bTransform.Position[0]),
		Y: int(bTransform.Position[1]),
	})

	return adjustedABounds.Overlaps(adjustedBBounds)
}

func (c *Collision) Update(ctx context.Context) error {
	ctx, task := trace.NewTask(ctx, "systems.Collision.Update")
	defer task.End()

	collisionCandidates := slices.Collect(ecs.Query2[components.CollisionComponent, components.Transform](ctx, c.EntityManager()))

	for i := range collisionCandidates {
		for j := i + 1; j < len(collisionCandidates); j++ {
			if c.checkCollision(ctx, collisionCandidates[i], collisionCandidates[j]) {

			}
		}
	}

	return nil
}
