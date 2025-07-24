package components

import (
	"golang.org/x/image/math/f64"
)

type Renderable struct{}

type Render struct {
	OnScreenPosition f64.Vec2
}

func (r *Render) Reset() {
	r.OnScreenPosition[0] = 0
	r.OnScreenPosition[1] = 0
}
