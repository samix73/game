package components

import "github.com/samix73/game/ecs"

func init() {
	ecs.RegisterComponent[Score]()
}

var _ ecs.Component = (*Score)(nil)

type Score struct {
	Distance float64
}

func (s *Score) Init() {
	s.Distance = 0
}

func (s *Score) Reset() {
	s.Distance = 0
}
