package systems

import (
	"fmt"
	"math/rand/v2"
	"slices"

	"github.com/jakecoffman/cp"
	"github.com/samix73/game/ecs"
	"github.com/samix73/game/game/components"
)

const (
	maxComingObstacles        = 6
	maxObstacleSpacing        = 550
	minObstacleSpacing        = 350
	topMaxObstacleHeight      = 6
	topMinObstacleHeight      = 4
	botMaxObstacleHeight      = 8
	botMinObstacleHeight      = 4
	obstacleBlockSize         = 64
	topObstacleSpawnChance    = 0.88
	bottomObstacleSpawnChance = 0.9
)

func init() {
	ecs.RegisterSystem(NewLevelGenSystem)
}

var _ ecs.System = (*LevelGenSystem)(nil)

type LevelGenSystem struct {
	*ecs.BaseSystem
}

func NewLevelGenSystem(priority int) *LevelGenSystem {
	return &LevelGenSystem{
		BaseSystem: ecs.NewBaseSystem(priority),
	}
}

func (l *LevelGenSystem) addObstacle(entityManager *ecs.EntityManager, position cp.Vector, height int) error {
	if _, err := entityManager.LoadEntity("Obstacle"); err != nil {
		return fmt.Errorf("error loading obstacle entity: %w", err)
	}

	// colors := []string{"red", "yellow", "blue", "green"}
	// if _, err := entities.NewObstacleEntity(
	// 	entityManager,
	// 	colors[rand.IntN(len(colors))],
	// 	height,
	// 	position,
	// ); err != nil {
	// 	return fmt.Errorf("error creating obstacle entity: %w", err)
	// }

	return nil
}

func (l *LevelGenSystem) obstacleHeight(min, max int) int {
	return rand.IntN(max-min) + min
}

func (l *LevelGenSystem) obstacleOffset(cameraBounds cp.BB, height int) float64 {
	return (cameraBounds.Dy() - float64(height*obstacleBlockSize)) * 0.5
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

	for entity := range ecs.Query[components.Obstacle](em) {
		obstacleTransform := ecs.MustGetComponent[components.Transform](em, entity)

		if obstacleTransform.Position.X < cameraLeft {
			em.Remove(entity)
		}

		if obstacleTransform.Position.X > playerTransform.Position.X {
			comingObstacles = append(comingObstacles, comingObstacle{
				id:        entity,
				transform: obstacleTransform,
			})
		}

		if len(comingObstacles) > maxComingObstacles {
			return nil
		}
	}

	var furthestDistance float64

	if len(comingObstacles) > 0 {
		furthest := slices.MaxFunc(comingObstacles, func(a, b comingObstacle) int {
			return int(a.transform.Position.X - b.transform.Position.X)
		})

		furthestDistance = furthest.transform.Position.X
	}

	spacing := rand.Float64()*(maxObstacleSpacing-minObstacleSpacing) + minObstacleSpacing
	xPosition := furthestDistance + spacing

	// Top obstacle
	if rand.Float64() < topObstacleSpawnChance {
		height := l.obstacleHeight(topMinObstacleHeight, topMaxObstacleHeight)

		if err := l.addObstacle(em,
			cp.Vector{
				X: xPosition,
				Y: l.obstacleOffset(cameraComponent.Bounds, height),
			},
			height,
		); err != nil {
			return fmt.Errorf("error adding obstacle: %w", err)
		}

	}

	// Bottom obstacle
	if rand.Float64() < bottomObstacleSpawnChance {
		height := l.obstacleHeight(botMinObstacleHeight, botMaxObstacleHeight)

		if err := l.addObstacle(em,
			cp.Vector{
				X: xPosition,
				Y: -l.obstacleOffset(cameraComponent.Bounds, height),
			},
			height,
		); err != nil {
			return fmt.Errorf("error adding obstacle: %w", err)
		}
	}

	return nil
}

func (l *LevelGenSystem) Teardown() {}
