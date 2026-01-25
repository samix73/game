package ecs_test

import (
	"slices"
	"testing"

	"github.com/jakecoffman/cp"
	"github.com/samix73/game/ecs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	ecs.RegisterComponent[TransformComponent]()
	ecs.RegisterComponent[CameraComponent]()
}

type TransformComponent struct {
	Position cp.Vector
	Rotation float64
}

func (t *TransformComponent) Init() {
	t.Position = cp.Vector{X: 0, Y: 0}
	t.Rotation = 0
}

func (t *TransformComponent) Reset() {
	t.Position = cp.Vector{X: 0, Y: 0}
	t.Rotation = 0
}

type CameraComponent struct {
	Zoom float64
}

func (c *CameraComponent) Init() {
	c.Zoom = 1.0
}

func (c *CameraComponent) Reset() {
	c.Zoom = 1.0
}

func NewPlayerEntity(tb testing.TB, em *ecs.EntityManager) ecs.EntityID {
	tb.Helper()

	entityID, err := em.NewEntity()
	require.NoError(tb, err)

	transform, err := ecs.AddComponent[TransformComponent](em, entityID)
	require.NotNil(tb, transform)
	require.NoError(tb, err)

	return entityID
}

func NewCameraEntity(tb testing.TB, em *ecs.EntityManager) ecs.EntityID {
	tb.Helper()

	entityID, err := em.NewEntity()
	require.NoError(tb, err)

	transform, err := ecs.AddComponent[TransformComponent](em, entityID)
	require.NotNil(tb, transform)
	require.NoError(tb, err)

	camera, err := ecs.AddComponent[CameraComponent](em, entityID)
	assert.NotNil(tb, camera)
	assert.NoError(tb, err)

	return entityID
}

func TestEntityCreation(t *testing.T) {
	em := ecs.NewEntityManager()

	player := NewPlayerEntity(t, em)
	assert.Equal(t, player, ecs.EntityID(1))
	camera := NewCameraEntity(t, em)
	assert.Equal(t, camera, ecs.EntityID(2))
	empty, err := em.NewEntity()
	require.NoError(t, err)
	assert.Equal(t, empty, ecs.EntityID(3))
}

func TestQuerySingleComponentFromMultiComponentEntity(t *testing.T) {
	em := ecs.NewEntityManager()

	// Create an entity with both Transform and Camera components
	entityID := NewCameraEntity(t, em)

	// Query for only Transform component (entity has both Transform and Camera)
	transformEntities := slices.Collect(ecs.Query[TransformComponent](em))

	// Entity should be found when querying for just Transform, even though it also has Camera
	assert.Contains(t, transformEntities, entityID, "Entity with Transform+Camera should be found when querying for just Transform")
	assert.Equal(t, 1, len(transformEntities), "Should find exactly 1 entity with Transform component")

	// Query for both components
	bothComponents := slices.Collect(ecs.Query2[TransformComponent, CameraComponent](em))
	assert.Contains(t, bothComponents, entityID, "Entity should be found when querying for both components")
	assert.Equal(t, 1, len(bothComponents), "Should find exactly 1 entity with both components")

	// Create an entity with only Transform
	playerID := NewPlayerEntity(t, em)

	// Query for Transform should now return both entities
	transformEntities = slices.Collect(ecs.Query[TransformComponent](em))
	assert.Contains(t, transformEntities, entityID, "Camera entity should still be in Transform query")
	assert.Contains(t, transformEntities, playerID, "Player entity should be in Transform query")
	assert.Equal(t, 2, len(transformEntities), "Should find 2 entities with Transform component")

	// Query for both components should only return the camera entity
	bothComponents = slices.Collect(ecs.Query2[TransformComponent, CameraComponent](em))
	assert.Contains(t, bothComponents, entityID, "Only camera entity has both components")
	assert.NotContains(t, bothComponents, playerID, "Player entity should not be in query for both components")
	assert.Equal(t, 1, len(bothComponents), "Should find exactly 1 entity with both components")
}

func BenchmarkGetComponent(b *testing.B) {
	em := ecs.NewEntityManager()

	// Create a set of entities with Transform components
	for range 500_000 {
		NewPlayerEntity(b, em)
	}

	for range 500_000 {
		NewCameraEntity(b, em)
	}

	for range 1000 {
		_, err := em.NewEntity()
		require.NoError(b, err)
	}

	entityIDs := slices.Collect(ecs.Query[TransformComponent](em))

	b.Run("GetComponent single entity", func(b *testing.B) {
		b.ResetTimer()
		for b.Loop() {
			ecs.GetComponent[TransformComponent](em, entityIDs[0])
		}
	})

	b.Run("GetComponent all entities", func(b *testing.B) {
		b.ResetTimer()
		for b.Loop() {
			for _, entityID := range entityIDs {
				ecs.GetComponent[TransformComponent](em, entityID)
			}
		}
	})
}

func BenchmarkQueryEntities(b *testing.B) {
	em := ecs.NewEntityManager()

	// Create a set of entities with Transform components
	for range 500_000 {
		NewPlayerEntity(b, em)
	}

	for range 500_000 {
		NewCameraEntity(b, em)
	}

	for range 1000 {
		_, err := em.NewEntity()
		require.NoError(b, err)
	}

	b.Run("Query", func(b *testing.B) {
		for b.Loop() {
			_ = slices.Collect(ecs.Query[TransformComponent](em))
		}
	})

	b.Run("Query2", func(b *testing.B) {
		for b.Loop() {
			_ = slices.Collect(ecs.Query2[TransformComponent, CameraComponent](em))
		}
	})
}
