package helpers

import (
	"golang.org/x/image/math/f64"
)

type AABB struct {
	Min, Max f64.Vec2
}

func (a *AABB) Reset() {
	a.Min[0] = 0
	a.Min[1] = 0
	a.Max[0] = 0
	a.Max[1] = 0
}

func (a *AABB) Center() f64.Vec2 {
	return f64.Vec2{
		a.Dx() * 0.5,
		a.Dy() * 0.5,
	}
}

func (a *AABB) Dy() float64 {
	return a.Max[1] - a.Min[1]
}

func (a *AABB) Dx() float64 {
	return a.Max[0] - a.Min[0]
}

func (a *AABB) Add(p f64.Vec2) AABB {
	return AABB{
		Min: f64.Vec2{
			a.Min[0] + p[0],
			a.Min[1] + p[1],
		},
		Max: f64.Vec2{
			a.Max[0] + p[0],
			a.Max[1] + p[1],
		},
	}
}

func (a *AABB) Overlaps(other AABB) bool {
	return a.Min[0] < other.Max[0] &&
		a.Max[0] > other.Min[0] &&
		a.Min[1] < other.Max[1] &&
		a.Max[1] > other.Min[1]
}

func (a *AABB) SetSize(width, height float64) {
	hw, hh := width*0.5, height*0.5
	a.Min[0], a.Min[1] = -hw, -hh
	a.Max[0], a.Max[1] = hw, hh
}

func (a *AABB) Offset(dx, dy float64) AABB {
	return AABB{
		Min: f64.Vec2{a.Min[0] + dx, a.Min[1] + dy},
		Max: f64.Vec2{a.Max[0] + dx, a.Max[1] + dy},
	}
}
