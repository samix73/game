package ecs_test

import (
	"slices"
	"testing"

	"github.com/samix73/game/ecs"
	"golang.org/x/image/math/f64"
)

type TransformComponent struct {
	Position f64.Vec2
	Rotation float64
}

func (t *TransformComponent) Init() {
	if t == nil {
		t = new(TransformComponent)
	}

	t.Position = f64.Vec2{0, 0}
	t.Rotation = 0
}

func (t *TransformComponent) Reset() {
	t.Position = f64.Vec2{0, 0}
	t.Rotation = 0
}

type CameraComponent struct {
	Zoom float64
}

func (c *CameraComponent) Init() {
	if c == nil {
		c = new(CameraComponent)
	}

	c.Zoom = 1.0
}

func (c *CameraComponent) Reset() {
	c.Zoom = 1.0
}

func NewPlayerEntity(tb testing.TB, em *ecs.EntityManager) ecs.EntityID {
	tb.Helper()

	entityID := em.NewEntity()

	ecs.AddComponent[*TransformComponent](em, entityID)

	return entityID
}

func NewCameraEntity(tb testing.TB, em *ecs.EntityManager) ecs.EntityID {
	tb.Helper()

	entityID := em.NewEntity()

	ecs.AddComponent[*TransformComponent](em, entityID)
	ecs.AddComponent[*CameraComponent](em, entityID)

	return entityID
}

func NewEmptyEntity(tb testing.TB, em *ecs.EntityManager) ecs.EntityID {
	tb.Helper()

	return em.NewEntity()
}

func BenchmarkEntityCreation(b *testing.B) {
	em := ecs.NewEntityManager()

	b.Run("Create Player Entity", func(b *testing.B) {
		for b.Loop() {
			NewPlayerEntity(b, em)
		}
	})

	b.Run("Create Camera Entity", func(b *testing.B) {
		for b.Loop() {
			NewCameraEntity(b, em)
		}
	})

	b.Run("Create Empty Entity", func(b *testing.B) {
		for b.Loop() {
			NewEmptyEntity(b, em)
		}
	})
}

func BenchmarkQueryEntities(b *testing.B) {
	em := ecs.NewEntityManager()

	// Create a set of entities with Transform components
	for range 1_000_000 {
		NewPlayerEntity(b, em)
	}

	for range 1000 {
		NewEmptyEntity(b, em)
	}

	b.Run("Query Only", func(b *testing.B) {
		for b.Loop() {
			for entityID := range ecs.Query[*TransformComponent](em) {
				_ = entityID // Just consume the entityID
			}
		}
	})

	b.Run("GetComponent Only", func(b *testing.B) {
		// Pre-collect entity IDs
		entityIDs := slices.Collect(ecs.Query[*TransformComponent](em))

		b.ResetTimer()
		for b.Loop() {
			for _, entityID := range entityIDs {
				ecs.GetComponent[*TransformComponent](em, entityID)
			}
		}
	})

	b.Run("Query + GetComponent", func(b *testing.B) {
		for b.Loop() {
			for entityID := range ecs.Query[*TransformComponent](em) {
				ecs.GetComponent[*TransformComponent](em, entityID)
			}
		}
	})
}
