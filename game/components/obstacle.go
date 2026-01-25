package components

import "github.com/samix73/game/ecs"

func init() {
	if err := ecs.RegisterComponent[Obstacle](); err != nil {
		panic(err)
	}
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
