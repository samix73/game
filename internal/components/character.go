package components

import (
	"errors"

	"github.com/hajimehoshi/ebiten/v2"
)

var _ Component = (*CharacterComponent)(nil)

type CharacterComponent struct {
	id       ComponentID
	Name     string
	Position *PositionComponent
}

type CharacterRepository struct {
	positionsRepository *PositionRepository
	characters          []CharacterComponent
}

func (c *CharacterRepository) Update() error {
	var joinedError error
	for _, character := range c.characters {
		if err := character.Update(); err != nil {
			joinedError = errors.Join(joinedError, err)
		}
	}

	return joinedError
}

func (c *CharacterRepository) Draw(screen *ebiten.Image) {
	for _, character := range c.characters {
		character.Draw(screen)
	}
}

func NewCharacterRepository(positionsRepository *PositionRepository) *CharacterRepository {
	return &CharacterRepository{
		positionsRepository: positionsRepository,
		characters:          make([]CharacterComponent, 0),
	}
}

func (c *CharacterComponent) ID() ComponentID {
	return c.id
}

func (c *CharacterComponent) Update() error { return nil }

func (c *CharacterComponent) Draw(screen *ebiten.Image) {}

func (r *CharacterRepository) New(name string) *CharacterComponent {
	character := CharacterComponent{
		id:   ComponentID(len(r.characters)),
		Name: name,
	}
	character.Position = r.positionsRepository.New(&character)

	r.characters = append(r.characters, character)

	return &character
}
