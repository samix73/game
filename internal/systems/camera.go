package systems

var _ ecs.System = (*Camera)(nil)

type Camera struct {
	*ecs.BaseSystem
}

func NewCameraSystem(priority int, entityManager *ecs.EntityManager) *Camera {
	return &Camera{
		BaseSystem: esc.NewBaseSystem(ecs.NextID(), priority, entityManager),
	}
}

func (c *Camera) Update() error {
	return nil
}

func (c *Camera) Teardown() {
}