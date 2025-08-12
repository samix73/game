package components

import "github.com/samix73/game/ecs"

var _ ecs.Component = (*Obstacle)(nil)

type Obstacle struct {
	Color  string
	Height int
}

func (o *Obstacle) Reset() {
	o.Color = ""
	o.Height = 0
}
