package components

import "golang.org/x/image/math/f64"

type Movable struct {
	Speed        float64
	Direction    f64.Vec2
	Acceleration f64.Vec2
}

func (m *Movable) Reset() {
	m.Speed = 0.0
	m.Direction[0] = 0
	m.Direction[1] = 0
	m.Acceleration[0] = 0
	m.Acceleration[1] = 0
}
