package ecs_test

import (
	"testing"

	"github.com/samix73/game/ecs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func highZoomFilter(c *CameraComponent) bool {
	return c.Zoom < 1.0
}

func lowZoomFilter(c *CameraComponent) bool {
	return c.Zoom > 0.5
}

func TestAnd(t *testing.T) {
	camera1 := &CameraComponent{Zoom: 1.5}
	camera2 := &CameraComponent{Zoom: 0.4}
	camera3 := &CameraComponent{Zoom: 0.6}

	assert.False(t, highZoomFilter(camera1))
	assert.True(t, highZoomFilter(camera2))
	assert.True(t, highZoomFilter(camera3))

	assert.True(t, lowZoomFilter(camera1))
	assert.False(t, lowZoomFilter(camera2))
	assert.True(t, lowZoomFilter(camera3))

	assert.False(t, ecs.And(highZoomFilter, lowZoomFilter)(camera1))
	assert.False(t, ecs.And(highZoomFilter, lowZoomFilter)(camera2))
	assert.True(t, ecs.And(highZoomFilter, lowZoomFilter)(camera3))
}

func TestOr(t *testing.T) {
	camera1 := &CameraComponent{Zoom: 1.5}
	camera2 := &CameraComponent{Zoom: 0.4}
	camera3 := &CameraComponent{Zoom: 0.6}

	assert.False(t, highZoomFilter(camera1))
	assert.True(t, highZoomFilter(camera2))
	assert.True(t, highZoomFilter(camera3))

	assert.True(t, lowZoomFilter(camera1))
	assert.False(t, lowZoomFilter(camera2))
	assert.True(t, lowZoomFilter(camera3))

	assert.True(t, ecs.Or(highZoomFilter, lowZoomFilter)(camera1))
	assert.True(t, ecs.Or(highZoomFilter, lowZoomFilter)(camera2))
	assert.True(t, ecs.Or(highZoomFilter, lowZoomFilter)(camera3))
}

func TestNot(t *testing.T) {
	camera1 := &CameraComponent{Zoom: 1.5}
	camera2 := &CameraComponent{Zoom: 0.4}
	camera3 := &CameraComponent{Zoom: 0.6}

	assert.True(t, ecs.Not(highZoomFilter)(camera1))
	assert.False(t, ecs.Not(highZoomFilter)(camera2))
	assert.False(t, ecs.Not(highZoomFilter)(camera3))

	assert.False(t, ecs.Not(lowZoomFilter)(camera1))
	assert.True(t, ecs.Not(lowZoomFilter)(camera2))
	assert.False(t, ecs.Not(lowZoomFilter)(camera3))
}

func TestWhere(t *testing.T) {
	em := ecs.NewEntityManager()

	camera1EntityID, err := em.NewEntity()
	require.NoError(t, err)
	camera1, err := ecs.AddComponent[CameraComponent](em, camera1EntityID)
	require.NoError(t, err)
	camera1.Zoom = 1.5

	camera2EntityID, err := em.NewEntity()
	camera2, err := ecs.AddComponent[CameraComponent](em, camera2EntityID)
	require.NoError(t, err)
	camera2.Zoom = 0.4

	camera3EntityID, err := em.NewEntity()
	camera3, err := ecs.AddComponent[CameraComponent](em, camera3EntityID)
	require.NoError(t, err)
	camera3.Zoom = 0.6

	cameras := ecs.Where(em, ecs.Query[CameraComponent](em), ecs.And(highZoomFilter, lowZoomFilter))

	gotCameras := make([]*CameraComponent, 0)
	for c := range cameras {
		gotCameras = append(gotCameras, ecs.MustGetComponent[CameraComponent](em, c))
	}

	assert.Len(t, gotCameras, 1)
	assert.Equal(t, camera3, gotCameras[0])
}
