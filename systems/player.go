package systems

import (
	"log/slog"

	"github.com/samix73/game/components"
	"github.com/samix73/game/ecs"
	"github.com/samix73/game/helpers"
	"github.com/samix73/game/keys"
	"golang.org/x/image/math/f64"
)

var _ ecs.System = (*Player)(nil)

type Player struct {
	*ecs.BaseSystem

	playerEntity        ecs.EntityID
	jumpForce           float64
	forwardAcceleration float64
	cameraOffset        f64.Vec2
	maxSpeed            float64
}

func NewPlayerSystem(priority int, entityManager *ecs.EntityManager,
	jumpForce float64, forwardAcceleration float64, cameraOffset f64.Vec2, maxSpeed float64) *Player {
	return &Player{
		BaseSystem:          ecs.NewBaseSystem(ecs.NextID(), priority, entityManager),
		jumpForce:           jumpForce,
		forwardAcceleration: forwardAcceleration * helpers.DeltaTime,
		cameraOffset:        cameraOffset,
		maxSpeed:            maxSpeed,
	}
}

func (p *Player) Teardown() {}

func (p *Player) getPlayerEntity() ecs.EntityID {
	if p.playerEntity == ecs.UndefinedID {
		playerEntity, ok := ecs.First(ecs.Query[components.Player](p.EntityManager()))
		if !ok {
			return ecs.UndefinedID
		}

		p.playerEntity = playerEntity
	}

	return p.playerEntity
}

func (p *Player) moveForward(rigidBody *components.RigidBody) {
	if rigidBody.Velocity[0] <= p.maxSpeed {
		rigidBody.ApplyAcceleration(f64.Vec2{p.forwardAcceleration, 0})
	}
}

func (p *Player) jump(rigidBody *components.RigidBody) {
	if keys.IsPressed(keys.PlayerJumpAction) {
		rigidBody.Velocity[1] = 0
		rigidBody.ApplyImpulse(f64.Vec2{0, p.jumpForce})
		slog.Debug("Jump!",
			slog.Any("velocity", rigidBody.Velocity),
		)
	}
}

func (p *Player) cameraFollow() {
	camera, ok := ecs.First(ecs.Query[components.ActiveCamera](p.EntityManager()))
	if !ok {
		return
	}

	playerTransform := ecs.MustGetComponent[components.Transform](p.EntityManager(), p.getPlayerEntity())
	cameraTransform := ecs.MustGetComponent[components.Transform](p.EntityManager(), camera)

	cameraTransform.SetPosition(playerTransform.Position[0]+p.cameraOffset[0], p.cameraOffset[1])
}

func (p *Player) Update() error {
	player := p.getPlayerEntity()
	if player == ecs.UndefinedID {
		return nil
	}

	rigidBody := ecs.MustGetComponent[components.RigidBody](p.EntityManager(), player)

	p.jump(rigidBody)
	p.moveForward(rigidBody)
	p.cameraFollow()

	return nil
}
