package components

import (
	"github.com/jakecoffman/cp"
	"github.com/samix73/game/ecs"
)

func init() {
	ecs.RegisterComponent[Collider]()
	ecs.RegisterComponent[Collision]()
}

var _ ecs.Component = (*Collider)(nil)

type Collider struct {
	Bounds cp.BB
}

func (c *Collider) Init() {
	c.Bounds = cp.BB{}
}

// SetSize sets the size of the collider bounds, centered at (0,0).
func (c *Collider) SetSize(width, height float64) {
	hw := width / 2
	hh := height / 2
	c.Bounds = cp.BB{L: -hw, B: -hh, R: hw, T: hh}
}

func (c *Collider) Reset() {
	c.Bounds = cp.BB{}
}

var _ ecs.Component = (*Collision)(nil)

type Collision struct {
	Entity      ecs.EntityID
	Penetration float64
	Normal      cp.Vector
}

func (c *Collision) Init() {}

func (c *Collision) Reset() {
	c.Entity = 0
	c.Penetration = 0
	c.Normal = cp.Vector{}
}
