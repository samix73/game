package systems

import (
	"testing"

	ecs "github.com/samix73/ebiten-ecs"
	"github.com/samix73/game/client/components"
	"github.com/stretchr/testify/require"
	"golang.org/x/image/math/f64"
)

func testRigidbodyEntity(t *testing.T, em *ecs.EntityManager, gravity bool) ecs.EntityID {
	t.Helper()

	entity := em.NewEntity()
	rigidBody := ecs.AddComponent[components.RigidBody](em, entity)
	rigidBody.Gravity = gravity

	return entity
}

func TestNewGravitySystem_CreatesInstance(t *testing.T) {
	const prio = 42
	g := NewGravitySystem(prio)
	require.NotNil(t, g)
	require.NotNil(t, g.BaseSystem)
	require.Equal(t, prio, g.Priority())

	require.Equal(t, gravity, g.dv)
}

func TestGravity_Update(t *testing.T) {
	em := ecs.NewEntityManager()

	gravityEnabledRigidBody := testRigidbodyEntity(t, em, true)
	gravityDisabledRigidBody := testRigidbodyEntity(t, em, false)

	game := ecs.NewGame(&ecs.GameConfig{})
	systemManager := ecs.NewSystemManager(em, game)

	g := NewGravitySystem(0)
	systemManager.Add(g)

	err := g.Update()
	require.NoError(t, err)

	expectedVelocity := f64.Vec2{
		gravity[0] * game.DeltaTime(),
		gravity[1] * game.DeltaTime(),
	}
	enableGravityRigidBody := ecs.MustGetComponent[components.RigidBody](em, gravityEnabledRigidBody)
	require.Equal(t, expectedVelocity, enableGravityRigidBody.Velocity)

	disableGravityRigidBody := ecs.MustGetComponent[components.RigidBody](em, gravityDisabledRigidBody)
	require.Equal(t, f64.Vec2{0, 0}, disableGravityRigidBody.Velocity)

}
