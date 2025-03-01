package components

import (
	"errors"

	"github.com/hajimehoshi/ebiten/v2"
)

type Repository struct {
	positionsRepo  *PositionRepository
	charactersRepo *CharacterRepository
}

func NewRepository() *Repository {
	return &Repository{}
}

func (l *Repository) Update() error {
	var joinedError error
	if err := l.Characters().Update(); err != nil {
		joinedError = errors.Join(joinedError, err)
	}

	return joinedError
}

func (l *Repository) Draw(screen *ebiten.Image) {
	l.Characters().Draw(screen)
}

func (r *Repository) Positions() *PositionRepository {
	if r.positionsRepo == nil {
		r.positionsRepo = NewPositionRepository()
	}

	return r.positionsRepo
}

func (r *Repository) Characters() *CharacterRepository {
	if r.charactersRepo == nil {
		r.charactersRepo = NewCharacterRepository(r.Positions())
	}

	return r.charactersRepo
}
