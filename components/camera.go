package components

type Camera struct {
	Zoom float64
}

func (c *Camera) Reset() {
	c.Zoom = 1.0
}

type ActiveCamera struct{}
