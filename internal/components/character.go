package components

var _ Component = (*Character)(nil)

type Character struct {
	id       ComponentID
	Name     string
	Position *PositionComponent
}

type CharacterRepository struct {
	positionsRepository *PositionRepository
	characters          []Character
}

func NewCharacterRepository(positionsRepository *PositionRepository) *CharacterRepository {
	return &CharacterRepository{
		positionsRepository: positionsRepository,
		characters:          make([]Character, 0),
	}
}

func (c *Character) ID() ComponentID {
	return c.id
}

func (r *CharacterRepository) New(name string) *Character {
	character := Character{
		id:   ComponentID(len(r.characters)),
		Name: name,
	}
	character.Position = r.positionsRepository.New(&character)

	r.characters = append(r.characters, character)

	return &character
}
