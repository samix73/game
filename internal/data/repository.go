package data

import (
	"github.com/samix73/game/internal/components"
)

type Repository struct {
	positionsRepo  *components.PositionRepository
	charactersRepo *components.CharacterRepository
}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) Positions() *components.PositionRepository {
	if r.positionsRepo == nil {
		r.positionsRepo = components.NewPositionRepository()
	}

	return r.positionsRepo
}

func (r *Repository) Characters() *components.CharacterRepository {
	if r.charactersRepo == nil {
		r.charactersRepo = components.NewCharacterRepository(r.Positions())
	}

	return r.charactersRepo
}
