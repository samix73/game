package components

import "github.com/samix73/game/ecs"

func init() {
	if err := ecs.RegisterComponent[Player](); err != nil {
		panic(err)
	}
}

var _ ecs.Component = (*Player)(nil)

type Player struct{}

func (p *Player) Init() {}

func (p *Player) Reset() {}
