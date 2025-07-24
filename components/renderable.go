package components

import (
	"github.com/samix73/game/ecs"
	"golang.org/x/image/math/f64"
)

var _ ecs.Component = (*Renderable)(nil)

type Renderable struct {
	CameraPosition f64.Vec2
}

func (r *Renderable) Reset() {
	r.CameraPosition[0] = 0
	r.CameraPosition[1] = 0
}
