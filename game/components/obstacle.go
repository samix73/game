package components

import ecs "github.com/samix73/ebiten-ecs"

func init() {
	ecs.RegisterComponent[Obstacle]()
}

var _ ecs.Component = (*Obstacle)(nil)

type Obstacle struct {
	Color  string
	Height int
}

func (o *Obstacle) Init() {}

func (o *Obstacle) Reset() {
	o.Color = ""
	o.Height = 0
}
