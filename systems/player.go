package systems

import (
	"context"
	"fmt"
	"log/slog"
	"runtime/trace"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/samix73/game/components"
	"github.com/samix73/game/ecs"
	"github.com/samix73/game/helpers"
	"golang.org/x/image/math/f64"
)

var _ ecs.System = (*Player)(nil)

type Player struct {
	*ecs.BaseSystem

	playerEntity        ecs.EntityID
	jumpKey             ebiten.Key
	jumpForce           float64
	forwardAcceleration float64
}

func NewPlayerSystem(ctx context.Context, priority int, entityManager *ecs.EntityManager,
	jumpKey ebiten.Key, jumpForce float64, forwardAcceleration float64) *Player {
	ctx, task := trace.NewTask(ctx, "systems.NewPlayerSystem")
	defer task.End()

	return &Player{
		BaseSystem:          ecs.NewBaseSystem(ctx, ecs.NextID(ctx), priority, entityManager),
		jumpKey:             jumpKey,
		jumpForce:           jumpForce,
		forwardAcceleration: forwardAcceleration,
	}
}

func (p *Player) Teardown() {

}

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

	rigidBody.ApplyAcceleration(f64.Vec2{p.forwardAcceleration, 0})
}

func (p *Player) applyInput(ctx context.Context, rigidBody *components.RigidBody) {
	region := trace.StartRegion(ctx, "systems.Player.applyInput")
	defer region.End()

	keys := inpututil.AppendJustPressedKeys([]ebiten.Key{})
	if slices.Contains(keys, p.jumpKey) {
		rigidBody.Velocity[1] *= 0.1 // Reset vertical velocity before applying jump force
		rigidBody.ApplyImpulse(f64.Vec2{0, -p.jumpForce})
		slog.Debug("Jump!",
			slog.String("velocity", fmt.Sprintf("(%f, %f)",
				rigidBody.Velocity[0], rigidBody.Velocity[1])),
		)
	}
}

func (p *Player) Update(ctx context.Context) error {
	ctx, task := trace.NewTask(ctx, "systems.Player.Update")
	defer task.End()

	player := p.getPlayerEntity(ctx)
	if player == ecs.UndefinedID {
		return nil
	}

	rigidBody := ecs.MustGetComponent[components.RigidBody](ctx, p.EntityManager(), player)

	p.applyInput(ctx, rigidBody)
	p.moveForward(ctx, rigidBody)

	return nil
}
