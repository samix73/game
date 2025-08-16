package helpers

import (
	"image"

	"golang.org/x/image/math/f64"
)

type AABB struct {
	Min, Max f64.Vec2
}

func (a *AABB) Add(p f64.Vec2) {
	a.Min[0] += p[0]
	a.Min[1] += p[1]
	a.Max[0] += p[0]
	a.Max[1] += p[1]
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
