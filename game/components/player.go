package components

import "github.com/samix73/game/ecs"

func init() {
	ecs.RegisterComponent[Player]()
}

var _ ecs.Component = (*Player)(nil)

type Player struct{}

func (p *Player) Init() {}

func (p *Player) Reset() {}
