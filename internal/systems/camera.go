package systems

var _ ecs.System = (*Camera)(nil)

type Camera struct {
	id ecs.SystemID
	priority int
}

func NewCameraSystem(priority int) *Camera {
	return &Camera{
		id: ecs.NextID(),
		priority: priority,
	}
}

func (c *Camera) ID() SystemID {
	return c.id
}

func (c *Camera) Priority() int {
	return c.priority
}

func (c *Camera) Update() error {
	return nil
}

func (c *Camera) Teardown() {

}