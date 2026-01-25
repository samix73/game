package systems

import (
	"fmt"
	"image/color"
	"math/rand/v2"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jakecoffman/cp"
	"github.com/samix73/game/ecs"
	"github.com/samix73/game/game/components"
)

const (
	maxComingObstacles = 4
	pipeSpacing        = 400.0 // Distance between pipe pairs
	pipeGap            = 250.0 // Vertical gap between top and bottom pipes
	pipeWidth          = 80.0
	pipeHeight         = 500.0
)

func init() {
	ecs.RegisterSystem(NewLevelGenSystem)
}

var _ ecs.System = (*LevelGenSystem)(nil)

type LevelGenSystem struct {
	*ecs.BaseSystem
	lastSpawnX float64
}

func NewLevelGenSystem(priority int) *LevelGenSystem {
	return &LevelGenSystem{
		BaseSystem: ecs.NewBaseSystem(priority),
		lastSpawnX: 0, // Will be set based on player position
	}
}

// createPipe creates a single pipe obstacle at the given position
func (l *LevelGenSystem) createPipe(em *ecs.EntityManager, position cp.Vector) error {
	entityID, err := em.NewEntity()
	if err != nil {
		return fmt.Errorf("error creating entity: %w", err)
	}

	// Add Transform
	transform, err := ecs.AddComponent[components.Transform](em, entityID)
	if err != nil {
		return fmt.Errorf("error adding transform: %w", err)
	}
	transform.SetPosition(position.X, position.Y)

	// Add Collider
	collider, err := ecs.AddComponent[components.Collider](em, entityID)
	if err != nil {
		return fmt.Errorf("error adding collider: %w", err)
	}
	collider.SetSize(pipeWidth, pipeHeight)

	// Add Obstacle tag
	ecs.AddComponent[components.Obstacle](em, entityID)

	// Add Renderable with a simple colored sprite
	renderable, err := ecs.AddComponent[components.Renderable](em, entityID)
	if err != nil {
		return fmt.Errorf("error adding renderable: %w", err)
	}

	// Create a simple green rectangle sprite for the pipe
	pipeImage := ebiten.NewImage(int(pipeWidth), int(pipeHeight))
	pipeImage.Fill(color.RGBA{34, 139, 34, 255}) // Green color
	renderable.Sprite = pipeImage
	renderable.Order = 1 // Render behind player

	return nil
}

func (l *LevelGenSystem) Update() error {
	em := l.EntityManager()

	camera, ok := ecs.First(ecs.Query[components.ActiveCamera](em))
	if !ok {
		return nil
	}

	cameraComponent := ecs.MustGetComponent[components.Camera](em, camera)
	cameraTransform := ecs.MustGetComponent[components.Transform](em, camera)

	player, ok := ecs.First(ecs.Query[components.Player](em))
	if !ok {
		return nil
	}

	playerTransform := ecs.MustGetComponent[components.Transform](em, player)

	cameraLeft := cameraTransform.Position.X - cameraComponent.Bounds.L

	type comingObstacle struct {
		id        ecs.EntityID
		transform *components.Transform
	}

	comingObstacles := make([]comingObstacle, 0)

	// Clean up obstacles that are off-screen and count upcoming ones
	for entity := range ecs.Query[components.Obstacle](em) {
		obstacleTransform := ecs.MustGetComponent[components.Transform](em, entity)

		if obstacleTransform.Position.X < cameraLeft-100 {
			if err := em.Remove(entity); err != nil {
				return fmt.Errorf("error removing obstacle: %w", err)
			}
			continue
		}

		if obstacleTransform.Position.X > playerTransform.Position.X {
			comingObstacles = append(comingObstacles, comingObstacle{
				id:        entity,
				transform: obstacleTransform,
			})
		}
	}

	// Determine the furthest obstacle position
	var furthestDistance float64 = l.lastSpawnX

	// Initialize spawn position based on player if this is the first time
	if l.lastSpawnX == 0 {
		furthestDistance = playerTransform.Position.X + 200 // Start spawning 200 units ahead
	}

	if len(comingObstacles) > 0 {
		furthest := slices.MaxFunc(comingObstacles, func(a, b comingObstacle) int {
			return int(a.transform.Position.X - b.transform.Position.X)
		})
		furthestDistance = furthest.transform.Position.X
	}

	// Spawn new pipe pairs if we have room
	if len(comingObstacles) < maxComingObstacles {
		xPosition := furthestDistance + pipeSpacing

		// Random gap center position (vertical)
		screenHeight := cameraComponent.Bounds.T - cameraComponent.Bounds.B
		gapCenterMin := -screenHeight/2 + pipeHeight/2 + pipeGap/2
		gapCenterMax := screenHeight/2 - pipeHeight/2 - pipeGap/2
		gapCenter := rand.Float64()*(gapCenterMax-gapCenterMin) + gapCenterMin

		// Top pipe (above the gap)
		topPipeY := gapCenter + pipeGap/2 + pipeHeight/2
		if err := l.createPipe(em, cp.Vector{X: xPosition, Y: topPipeY}); err != nil {
			return fmt.Errorf("error creating top pipe: %w", err)
		}

		// Bottom pipe (below the gap)
		bottomPipeY := gapCenter - pipeGap/2 - pipeHeight/2
		if err := l.createPipe(em, cp.Vector{X: xPosition, Y: bottomPipeY}); err != nil {
			return fmt.Errorf("error creating bottom pipe: %w", err)
		}

		l.lastSpawnX = xPosition
	}

	return nil
}

func (l *LevelGenSystem) Teardown() {}
