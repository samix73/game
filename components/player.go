package components

import ecs "github.com/samix73/ebiten-ecs"

var _ ecs.Component = (*Player)(nil)

type Player struct{}

func (p *Player) Reset() {}
