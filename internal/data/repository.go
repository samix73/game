package data

type Repository struct {
	entities map[EntityType][]Entity
}

func NewRepository() *Repository {
	return &Repository{
		entities: make(map[EntityType][]Entity),
	}
}

func (r *Repository) Add(e Entity) Entity {
	e.SetID(EntityID(len(r.entities[e.Type()])))
	r.entities[e.Type()] = append(r.entities[e.Type()], e)

	return e
}
