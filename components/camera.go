package components

// Camera represents the viewable area of the game world.
type Camera struct {
	Zoom float64
}

func (c *Camera) Reset() {
	c.Zoom = 1.0
}

type ActiveCamera struct{}
