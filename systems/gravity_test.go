package systems

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/image/math/f64"
)

func TestNewGravitySystem_CreatesInstance(t *testing.T) {
	const prio = 42
	g := NewGravitySystem(prio)
	require.NotNil(t, g)
	require.NotNil(t, g.BaseSystem)
	require.Equal(t, prio, g.Priority())

	expected := f64.Vec2{0, -981}
	require.Equal(t, expected, g.dv)
}
