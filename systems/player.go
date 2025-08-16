package systems

import (
	"context"
	"log/slog"
	"runtime/trace"

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

func NewPlayerSystem(ctx context.Context, priority int, entityManager *ecs.EntityManager,
	jumpForce float64, forwardAcceleration float64, cameraOffset f64.Vec2, maxSpeed float64) *Player {
	ctx, task := trace.NewTask(ctx, "systems.NewPlayerSystem")
	defer task.End()

	return &Player{
		BaseSystem:          ecs.NewBaseSystem(ctx, ecs.NextID(ctx), priority, entityManager),
		jumpForce:           jumpForce,
		forwardAcceleration: forwardAcceleration * helpers.DeltaTime,
		cameraOffset:        cameraOffset,
		maxSpeed:            maxSpeed,
	}
}

func (p *Player) Teardown() {}

func (p *Player) getPlayerEntity(ctx context.Context) ecs.EntityID {
	ctx, task := trace.NewTask(ctx, "systems.Player.getPlayerEntity")
	defer task.End()

	if p.playerEntity == ecs.UndefinedID {
		playerEntity, ok := helpers.First(ecs.Query[components.Player](ctx, p.EntityManager()))
		if !ok {
			return ecs.UndefinedID
		}

		p.playerEntity = playerEntity
	}

	return p.playerEntity
}

func (p *Player) moveForward(ctx context.Context, rigidBody *components.RigidBody) {
	region := trace.StartRegion(ctx, "systems.Player.moveForward")
	defer region.End()

	if rigidBody.Velocity[0] <= p.maxSpeed {
		rigidBody.ApplyAcceleration(f64.Vec2{p.forwardAcceleration, 0})
	}
}

func (p *Player) jump(ctx context.Context, rigidBody *components.RigidBody) {
	region := trace.StartRegion(ctx, "systems.Player.jump")
	defer region.End()

	if keys.IsPressed(keys.PlayerJumpAction) {
		rigidBody.Velocity[1] *= 0.1
		rigidBody.ApplyImpulse(f64.Vec2{0, p.jumpForce})
		slog.Debug("Jump!",
			slog.Any("velocity", rigidBody.Velocity),
		)
	}
}

func (p *Player) cameraFollow(ctx context.Context) {
	region := trace.StartRegion(ctx, "systems.Player.cameraFollow")
	defer region.End()

	camera, ok := helpers.First(ecs.Query[components.ActiveCamera](ctx, p.EntityManager()))
	if !ok {
		return
	}

	playerTransform := ecs.MustGetComponent[components.Transform](ctx, p.EntityManager(), p.getPlayerEntity(ctx))
	cameraTransform := ecs.MustGetComponent[components.Transform](ctx, p.EntityManager(), camera)

	cameraTransform.SetPosition(playerTransform.Position[0]+p.cameraOffset[0], p.cameraOffset[1])
}

func (p *Player) Update(ctx context.Context) error {
	ctx, task := trace.NewTask(ctx, "systems.Player.Update")
	defer task.End()

	player := p.getPlayerEntity(ctx)
	if player == ecs.UndefinedID {
		return nil
	}

	rigidBody := ecs.MustGetComponent[components.RigidBody](ctx, p.EntityManager(), player)

	p.jump(ctx, rigidBody)
	p.moveForward(ctx, rigidBody)
	p.cameraFollow(ctx)

	return nil
}
