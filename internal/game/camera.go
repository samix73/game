package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/samix73/game/ecs"
	"golang.org/x/image/math/f64"
)

type Camera struct {
	Vec f64.Vec2

	// Width and Height represent the viewport dimensions
	Width, Height int

	// Zoom level of the camera (1.0 = normal size)
	Zoom float64

	// Associated world that this camera renders
	world *ecs.World
}

// NewCamera creates a new camera attached to the given world with default settings
func NewCamera(world *ecs.World, width, height int) *Camera {
	return &Camera{
		Vec:    f64.Vec2{0, 0},
		Width:  width,
		Height: height,
		Zoom:   1.0,
		world:  world,
	}
}

func (c *Camera) Update() error {
	return nil
}

func (c *Camera) Draw(screen *ebiten.Image) {
}
