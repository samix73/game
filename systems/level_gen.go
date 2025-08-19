package systems

import (
	"fmt"
	"math/rand/v2"
	"slices"

	"github.com/samix73/game/components"
	"github.com/samix73/game/ecs"
	"github.com/samix73/game/entities"
	"github.com/samix73/game/helpers"
	"golang.org/x/image/math/f64"
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

var _ ecs.System = (*LevelGen)(nil)

type LevelGen struct {
	*ecs.BaseSystem
}

func NewLevelGenSystem(priority int, entityManager *ecs.EntityManager) *LevelGen {
	return &LevelGen{
		BaseSystem: ecs.NewBaseSystem(ecs.NextID(), priority, entityManager),
	}
}

func (l *LevelGen) addObstacle(entityManager *ecs.EntityManager, position f64.Vec2, height int) error {
	colors := []string{"red", "yellow", "blue", "green"}
	if _, err := entities.NewObstacleEntity(
		entityManager,
		colors[rand.IntN(len(colors))],
		height,
		position,
	); err != nil {
		return fmt.Errorf("error creating obstacle entity: %w", err)
	}

	return nil
}

func (l *LevelGen) obstacleHeight(min, max int) int {
	return rand.IntN(max-min) + min
}

func (l *LevelGen) obstacleOffset(cameraBounds helpers.AABB, height int) float64 {
	return (cameraBounds.Dy() - float64(height*obstacleBlockSize)) * 0.5
}

func (l *LevelGen) Update() error {
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

	cameraLeft := cameraTransform.Position[0] - cameraComponent.Bounds.Dx()/2

	type comingObstacle struct {
		id        ecs.EntityID
		transform *components.Transform
	}

	comingObstacles := make([]comingObstacle, 0)

	for entity := range ecs.Query[components.Obstacle](em) {
		obstacleTransform := ecs.MustGetComponent[components.Transform](em, entity)

		if obstacleTransform.Position[0] < cameraLeft {
			em.Remove(entity)
		}

		if obstacleTransform.Position[0] > playerTransform.Position[0] {
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
			return int(a.transform.Position[0] - b.transform.Position[0])
		})

		furthestDistance = furthest.transform.Position[0]
	}

	spacing := rand.Float64()*(maxObstacleSpacing-minObstacleSpacing) + minObstacleSpacing
	xPosition := furthestDistance + spacing

	// Top obstacle
	if rand.Float64() < topObstacleSpawnChance {
		height := l.obstacleHeight(topMinObstacleHeight, topMaxObstacleHeight)

		if err := l.addObstacle(em,
			f64.Vec2{
				xPosition,
				l.obstacleOffset(cameraComponent.Bounds, height),
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
			f64.Vec2{
				xPosition,
				-l.obstacleOffset(cameraComponent.Bounds, height),
			},
			height,
		); err != nil {
			return fmt.Errorf("error adding obstacle: %w", err)
		}
	}

	return nil
}

func (l *LevelGen) Teardown() {}
