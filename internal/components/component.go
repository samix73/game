package components

import "github.com/hajimehoshi/ebiten/v2"

type ComponentID uint64

type Component interface {
	ID() ComponentID
	Update() error
	Draw(screen *ebiten.Image)
}
