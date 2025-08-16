package keys

import (
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Action struct {
	Keys        []ebiten.Key
	MouseButton []ebiten.MouseButton
}

var (
	PlayerJumpAction = Action{
		Keys: []ebiten.Key{
			ebiten.KeySpace,
			ebiten.KeyArrowUp,
			ebiten.KeyW,
		},
		MouseButton: []ebiten.MouseButton{ebiten.MouseButtonLeft},
	}
	PauseAction = Action{
		Keys:        []ebiten.Key{ebiten.KeyP, ebiten.KeyEscape},
		MouseButton: []ebiten.MouseButton{},
	}
)

func IsPressed(action Action) bool {
	return slices.ContainsFunc(action.Keys, inpututil.IsKeyJustPressed) ||
		slices.ContainsFunc(action.MouseButton, inpututil.IsMouseButtonJustPressed)
}
