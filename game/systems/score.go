package systems

import (
	ecs "github.com/samix73/ebiten-ecs"
	"github.com/samix73/game/game/components"
	"github.com/samix73/game/helpers"
)

var _ ecs.System = (*ScoreSystem)(nil)

func init() {
	ecs.RegisterSystem(NewScoreSystem)
}

type ScoreSystem struct {
	*ecs.BaseSystem
	lastCameraX float64
}

func NewScoreSystem(priority int) *ScoreSystem {
	return &ScoreSystem{
		BaseSystem:  ecs.NewBaseSystem(priority),
		lastCameraX: 0,
	}
}

func (s *ScoreSystem) Update() error {
	em := s.EntityManager()

	// Get the active camera to track distance
	camera, ok := helpers.First(ecs.Query[components.ActiveCamera](em))
	if !ok {
		return nil
	}

	cameraTransform := ecs.MustGetComponent[components.Transform](em, camera)

	// Calculate distance traveled since last frame
	distance := cameraTransform.Position.X - s.lastCameraX
	s.lastCameraX = cameraTransform.Position.X

	// Update player score
	for _, entity := range ecs.Query[components.Score](em) {
		score := ecs.MustGetComponent[components.Score](em, entity)
		if distance > 0 {
			score.Distance += distance
		}
	}

	return nil
}

func (s *ScoreSystem) Start() error {
	return nil
}

func (s *ScoreSystem) Teardown() {
}
