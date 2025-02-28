package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/samix73/game/internal/data"
)

type Level struct {
	repository *data.Repository
}

func (l *Level) Load() error {
	return nil
}

func (l *Level) Draw(screen *ebiten.Image) {
}

func (l *Level) Update() error {
	return nil
}
