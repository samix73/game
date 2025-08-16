package systems

import (
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

func NewCollisionSystem(priority int, entityManager *ecs.EntityManager) *Collision {
	return &Collision{
		BaseSystem: ecs.NewBaseSystem(ecs.NextID(), priority, entityManager),
	}
}

func (c *Collision) Teardown() {}

func (c *Collision) checkCollision(a, b collisionCandidate) bool {
	return a.bounds.Overlaps(b.bounds)
}

func (c *Collision) Update() error {
	em := c.EntityManager()

	active := make([]collisionCandidate, 0, 16)
	static := make([]collisionCandidate, 0, 1024)

	for entity := range ecs.Query2[components.Collider, components.Transform](em) {
		transform := ecs.MustGetComponent[components.Transform](em, entity)
		col := ecs.MustGetComponent[components.Collider](em, entity)

		adjustedBounds := col.Bounds.Add(transform.Position)

		if ecs.HasComponent[components.RigidBody](em, entity) {
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
			if c.checkCollision(a, b) {
				aCol := ecs.AddComponent[components.Collision](em, a.id)
				aCol.Enitity = b.id

				bCol := ecs.AddComponent[components.Collision](em, a.id)
				bCol.Enitity = a.id
			} else {
				ecs.RemoveComponent[components.Collision](em, a.id)
				ecs.RemoveComponent[components.Collision](em, b.id)
			}
		}
	}

	if len(active) < 2 {
		return nil
	}

	// Active vs Active
	for i := 0; i < len(active); i++ {
		for j := i + 1; j < len(active); j++ {
			if c.checkCollision(active[i], active[j]) {
				aCol := ecs.AddComponent[components.Collision](em, active[i].id)
				aCol.Enitity = active[j].id

				bCol := ecs.AddComponent[components.Collision](em, active[j].id)
				bCol.Enitity = active[i].id
			} else {
				ecs.RemoveComponent[components.Collision](em, active[i].id)
				ecs.RemoveComponent[components.Collision](em, active[j].id)
			}
		}
	}

	return nil
}
