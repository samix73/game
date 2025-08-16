package keys

import (
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Action []ebiten.Key

var (
	PlayerJumpAction = Action{ebiten.KeySpace, ebiten.KeyArrowUp}
	PauseAction      = Action{ebiten.KeyP}
)

func IsPressed(action Action) bool {
	return slices.ContainsFunc(action, inpututil.IsKeyJustPressed)
}
