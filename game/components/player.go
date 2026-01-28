package components

import ecs "github.com/samix73/ebiten-ecs"

func init() {
	ecs.RegisterComponent[Player]()
}

var _ ecs.Component = (*Player)(nil)

type Player struct{}

func (p *Player) Init() {}

func (p *Player) Reset() {}
