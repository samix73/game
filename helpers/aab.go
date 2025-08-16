package helpers

import (
	"image"

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

func (a *AABB) SetImageBounds(bounds image.Rectangle) {
	a.Min = f64.Vec2{float64(bounds.Min.X), float64(bounds.Min.Y)}
	a.Max = f64.Vec2{float64(bounds.Max.X), float64(bounds.Max.Y)}
}

func (a *AABB) Offset(dx, dy float64) AABB {
	return AABB{
		Min: f64.Vec2{a.Min[0] + dx, a.Min[1] + dy},
		Max: f64.Vec2{a.Max[0] + dx, a.Max[1] + dy},
	}
}
