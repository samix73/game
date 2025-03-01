package components

import (
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/math/f64"
)

var _ Component = (*PositionComponent)(nil)

type PositionComponent struct {
	id    ComponentID
	Vec2  f64.Vec2
	Owner Component
}

type PositionRepository struct {
	positions []PositionComponent
}

func NewPositionRepository() *PositionRepository {
	return &PositionRepository{
		positions: make([]PositionComponent, 0),
	}
}

func (c *PositionComponent) ID() ComponentID {
	return c.id
}

func (c *PositionComponent) Update() error {
	panic("don't call Update on PositionComponent")
}

func (c *PositionComponent) Draw(screen *ebiten.Image) {
	panic("don't call Draw on PositionComponent")
}

func (r *PositionRepository) New(owner Component) *PositionComponent {
	position := PositionComponent{
		id:    ComponentID(len(r.positions)),
		Vec2:  f64.Vec2{0, 0},
		Owner: owner,
	}
	r.positions = append(r.positions, position)

	return &position
}
