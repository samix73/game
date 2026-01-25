package systems

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/samix73/game/ecs"
	"github.com/samix73/game/game/components"
)

var _ ecs.DrawableSystem = (*GameOverSystem)(nil)

func init() {
	ecs.RegisterSystem(NewGameOverSystem)
}

type GameOverSystem struct {
	*ecs.BaseSystem
	gameOver bool
	score    float64
}

func NewGameOverSystem(priority int) *GameOverSystem {
	return &GameOverSystem{
		BaseSystem: ecs.NewBaseSystem(priority),
		gameOver:   false,
		score:      0,
	}
}

func (g *GameOverSystem) Update() error {
	// Skip if already game over
	if g.gameOver {
		return nil
	}

	em := g.EntityManager()

	// Check if player has collision with an obstacle
	for entity := range ecs.Query2[components.Player, components.Collision](em) {
		collision := ecs.MustGetComponent[components.Collision](em, entity)

		// Check if the collision is with an obstacle
		if ecs.HasComponent[components.Obstacle](em, collision.Entity) {
			// Game over!
			g.gameOver = true

			// Pause the game
			g.Game().SetTimeScale(0)

			// Get the score
			score, hasScore := ecs.GetComponent[components.Score](em, entity)
			if hasScore {
				g.score = score.Distance
			}

			break
		}
	}

	return nil
}

func (g *GameOverSystem) Draw(screen *ebiten.Image) {
	if !g.gameOver {
		return
	}

	cfg := g.Game().Config()
	screenWidth := cfg.ScreenWidth
	screenHeight := cfg.ScreenHeight

	// Draw game over text using ebitenutil
	gameOverText := "GAME OVER"
	scoreText := fmt.Sprintf("Score: %.0f", g.score)
	restartText := "Press R to Restart"

	// Draw centered text
	ebitenutil.DebugPrintAt(screen, gameOverText, screenWidth/2-80, screenHeight/2-100)
	ebitenutil.DebugPrintAt(screen, scoreText, screenWidth/2-60, screenHeight/2-40)
	ebitenutil.DebugPrintAt(screen, restartText, screenWidth/2-100, screenHeight/2+20)
}

func (g *GameOverSystem) Teardown() {
}
