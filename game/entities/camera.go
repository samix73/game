package entities

import (
	"fmt"

	"github.com/jakecoffman/cp"
	"github.com/samix73/game/ecs"
	"github.com/samix73/game/game/components"
)

func NewCameraEntity(em *ecs.EntityManager, active bool, cameraWidth, cameraHeight float64) (ecs.EntityID, error) {
	entityID, err := em.NewEntity()
	if err != nil {
		return 0, fmt.Errorf("error creating entity: %w", err)
	}
	if active {
		if _, err := ecs.AddComponent[components.ActiveCamera](em, entityID); err != nil {
			return entityID, fmt.Errorf("error adding active camera: %w", err)
		}
	}

	if _, err := ecs.AddComponent[components.Transform](em, entityID); err != nil {
		return entityID, fmt.Errorf("error adding transform: %w", err)
	}

	camera, err := ecs.AddComponent[components.Camera](em, entityID)
	if err != nil {
		return entityID, fmt.Errorf("error adding camera: %w", err)
	}
	camera.Zoom = 1.0
	camera.Bounds = cp.BB{L: 0, B: 0, R: cameraWidth, T: cameraHeight}

	return entityID, nil
}
