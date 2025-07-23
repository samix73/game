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

func (c *Camera) createDefualtCamera() ecs.EntityID {
	em := c.EntityManager()

	// Create a new camera entity
	camera := ecs.NewEntity(em)
	ecs.AddComponent[*components.Camera](em, camera)
	ecs.AddComponent[*components.ActiveCamera](em, camera)
	ecs.AddComponent[*components.Transform](em, camera)

	// Set a default position for the camera
	transform := em.GetComponent[*components.Transform](camera)
	transform.Position = ecs.Vector3{X: 0, Y: 0, Z: 10}

	return camera
}

func (c *Camera) activeCamera() ecs.EntityID {
	em := c.EntityManager()

	var activeCamera ecs.EntityID
	for camera := range ecs.Query2[*components.Camera, *components.ActiveCamera](em) {
		activeCamera = camera
		break
	}

	if activeCamera == ecs.UndefinedID {
		for camera := range ecs.Query[*components.Camera](em) {
			ecs.AddComponent[*components.ActiveCamera](em, camera)
			activeCamera = camera

			// We only want to set the first camera as active
			break
		}
	}

	// If no camera is found, we create a default one
	if activeCamera == ecs.UndefinedID {
		activeCamera = c.createDefualtCamera()
	}

	return activeCamera
}

func (c *Camera) Update() error {
	activeCamera := c.activeCamera()
	
	return nil
}

func (c *Camera) Teardown() {
}