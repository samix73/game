package components

import ecs "github.com/samix73/ebiten-ecs"

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
