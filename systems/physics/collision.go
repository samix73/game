package physics

import (
	"github.com/jakecoffman/cp"
	ecs "github.com/samix73/ebiten-ecs"
	"github.com/samix73/game/components"
)

var _ ecs.System = (*Collision)(nil)

type collisionCandidate struct {
	id     ecs.EntityID
	bounds cp.BB
}

type Collision struct {
	*ecs.BaseSystem
}

func NewCollisionSystem(priority int) *Collision {
	return &Collision{
		BaseSystem: ecs.NewBaseSystem(ecs.NextID(), priority),
	}
}

func (c *Collision) checkCollision(a, b collisionCandidate) bool {
	return a.bounds.Intersects(b.bounds)
}

func (c *Collision) registerCollision(a, b ecs.EntityID) {
	aCol := ecs.AddComponent[components.Collision](c.EntityManager(), a)
	aCol.Entity = b

	bCol := ecs.AddComponent[components.Collision](c.EntityManager(), b)
	bCol.Entity = a
}

func (c *Collision) removeCollision(a, b ecs.EntityID) {
	aCol, ok := ecs.GetComponent[components.Collision](c.EntityManager(), a)
	if ok && aCol.Entity == b {
		ecs.RemoveComponent[components.Collision](c.EntityManager(), a)
	}

	bCol, ok := ecs.GetComponent[components.Collision](c.EntityManager(), b)
	if ok && bCol.Entity == a {
		ecs.RemoveComponent[components.Collision](c.EntityManager(), b)
	}
}

func (c *Collision) Update() error {
	em := c.EntityManager()

	active := make([]collisionCandidate, 0, 16)
	static := make([]collisionCandidate, 0, 1024)

	for entity := range ecs.Query2[components.Collider, components.Transform](em) {
		transform := ecs.MustGetComponent[components.Transform](em, entity)
		col := ecs.MustGetComponent[components.Collider](em, entity)

		translatedBounds := col.Bounds.Offset(transform.Position)

		if ecs.HasComponent[components.RigidBody](em, entity) {
			active = append(active, collisionCandidate{
				id:     entity,
				bounds: translatedBounds,
			})
		} else {
			static = append(static, collisionCandidate{
				id:     entity,
				bounds: translatedBounds,
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
				c.registerCollision(a.id, b.id)
			} else {
				c.removeCollision(a.id, b.id)
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
				c.registerCollision(active[i].id, active[j].id)
			} else {
				c.removeCollision(active[i].id, active[j].id)
			}
		}
	}

	return nil
}
