package worlds

import (
	"bytes"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/images"
	"github.com/jakecoffman/cp"
	ecs "github.com/samix73/ebiten-ecs"
	"github.com/samix73/game/components"
	"github.com/samix73/game/entities"
	"github.com/samix73/game/systems"
	"github.com/samix73/game/systems/physics"
)

var _ ecs.World = (*MainWorld)(nil)

var tilesData = [][]int{
	{
		243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
		243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
		243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
		243, 218, 243, 243, 243, 243, 243, 243, 243, 243, 243, 218, 243, 244, 243,
		243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,

		243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
		243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
		243, 243, 244, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
		243, 243, 243, 243, 243, 243, 243, 243, 243, 219, 243, 243, 243, 219, 243,
		243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,

		243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
		243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
		243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
		243, 218, 243, 243, 243, 243, 243, 243, 243, 243, 243, 244, 243, 243, 243,
		243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 26, 27, 28, 29, 30, 31, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 51, 52, 53, 54, 55, 56, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 76, 77, 78, 79, 80, 81, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 101, 102, 103, 104, 105, 106, 0, 0, 0, 0,

		0, 0, 0, 0, 0, 126, 127, 128, 129, 130, 131, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 303, 303, 245, 242, 303, 303, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,

		0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
	},
}

type MainWorld struct {
	*ecs.BaseWorld
}

func (m *MainWorld) Init(g *ecs.Game) error {
	g.SetTimeScale(1)

	entityManager := ecs.NewEntityManager()
	systemManager := ecs.NewSystemManager(entityManager, g)

	m.BaseWorld = ecs.NewBaseWorld(entityManager, systemManager)

	m.registerSystems()

	m.addEntities()

	return nil
}

func (m *MainWorld) addEntities() {
	em := m.EntityManager()

	img, _, _ := image.Decode(bytes.NewReader(images.Tiles_png))
	atlas := ebiten.NewImageFromImage(img)

	const (
		screenWidth  = 15
		screenHeight = 15
	)

	for layer, tiles := range tilesData {
		entities.NewTileMapEntity(em, atlas, layer, screenWidth, screenHeight, tiles)
	}

	// Add test entities for collision testing
	m.addTestEntities()
}

func (m *MainWorld) addTestEntities() {
	em := m.EntityManager()

	// Test 1: Two rigidbodies moving towards each other
	entity1 := em.NewEntity()
	transform1 := ecs.AddComponent[components.Transform](em, entity1)
	transform1.SetPosition(100, 200)

	rb1 := ecs.AddComponent[components.RigidBody](em, entity1)
	rb1.Mass = 1.0
	rb1.Velocity = cp.Vector{X: 100, Y: 0} // Moving right
	rb1.Gravity = false

	collider1 := ecs.AddComponent[components.Collider](em, entity1)
	collider1.SetSize(32, 32)

	renderable := ecs.AddComponent[components.Renderable](em, entity1)
	renderable.Sprite = ebiten.NewImage(32, 32)
	renderable.Sprite.Fill(color.RGBA{255, 0, 0, 255}) // Red color

	entity2 := em.NewEntity()
	transform2 := ecs.AddComponent[components.Transform](em, entity2)
	transform2.SetPosition(300, 200)

	rb2 := ecs.AddComponent[components.RigidBody](em, entity2)
	rb2.Mass = 2.0
	rb2.Velocity = cp.Vector{X: -50, Y: 0} // Moving left
	rb2.Gravity = false

	collider2 := ecs.AddComponent[components.Collider](em, entity2)
	collider2.SetSize(32, 32)

	renderable2 := ecs.AddComponent[components.Renderable](em, entity2)
	renderable2.Sprite = ebiten.NewImage(32, 32)
	renderable2.Sprite.Fill(color.RGBA{0, 0, 255, 255}) // Blue color

	// Test 2: One rigidbody moving towards a static object
	entity3 := em.NewEntity()
	transform3 := ecs.AddComponent[components.Transform](em, entity3)
	transform3.SetPosition(100, 300)

	rb3 := ecs.AddComponent[components.RigidBody](em, entity3)
	rb3.Mass = 1.0
	rb3.Velocity = cp.Vector{X: 80, Y: 0} // Moving right
	rb3.Gravity = false

	collider3 := ecs.AddComponent[components.Collider](em, entity3)
	collider3.SetSize(32, 32)

	renderable3 := ecs.AddComponent[components.Renderable](em, entity3)
	renderable3.Sprite = ebiten.NewImage(32, 32)
	renderable3.Sprite.Fill(color.RGBA{0, 255, 0, 255}) // Green color

	// Static object (no rigidbody)
	entity4 := em.NewEntity()
	transform4 := ecs.AddComponent[components.Transform](em, entity4)
	transform4.SetPosition(300, 300)

	collider4 := ecs.AddComponent[components.Collider](em, entity4)
	collider4.SetSize(32, 32)

	renderable4 := ecs.AddComponent[components.Renderable](em, entity4)
	renderable4.Sprite = ebiten.NewImage(32, 32)
	renderable4.Sprite.Fill(color.RGBA{255, 255, 0, 255}) // Yellow color
}

func (m *MainWorld) registerSystems() {
	m.SystemManager().Add(
		systems.NewPauseSystem(0),
		physics.NewGravitySystem(1),
		physics.NewPhysicsSystem(2),
		physics.NewCollisionSystem(3),
		physics.NewCollisionResolverSystem(4),
		systems.NewCameraSystem(5),
		systems.NewTileSystem(6),
	)
}
