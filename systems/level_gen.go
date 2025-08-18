package systems

import (
	"fmt"

	"github.com/samix73/game/components"
	"github.com/samix73/game/ecs"
	"github.com/samix73/game/helpers"
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

func (l *LevelGen) addObstacle() {
	// colors := []string{"red", "yellow", "blue"}
	// for i := range 10_000 {
	// 	if _, err := entities.NewObstacleEntity(
	// 		entityManager,
	// 		colors[rand.IntN(len(colors))],
	// 		rand.IntN(8)+3,
	// 		f64.Vec2{float64(i * 450), 0},
	// 	); err != nil {
	// 		return nil, fmt.Errorf("error creating obstacle entity: %w", err)
	// 	}
	// }
}

func (l *LevelGen) Update() error {
	em := l.EntityManager()

	player, ok := helpers.First(ecs.Query[components.ActiveCamera](em))
	if !ok {
		return nil
	}

	cameraComponent := ecs.MustGetComponent[components.Camera](em, player)
	cameraTransform := ecs.MustGetComponent[components.Transform](em, player)

	fmt.Println(cameraComponent, cameraTransform)

	return nil
}

func (l *LevelGen) Teardown() {}
