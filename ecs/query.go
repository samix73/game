package ecs

type Query struct {
	world            *World
	componentTypeIDs []ComponentTypeID
}

func NewQuery(world *World) *Query {
	return &Query{
		world:            world,
		componentTypeIDs: make([]ComponentTypeID, 0),
	}
}

func (q *Query) WithComponentTypeIDs(componentTypeIDs ...ComponentTypeID) *Query {
	q.componentTypeIDs = append(q.componentTypeIDs, componentTypeIDs...)
	return q
}

func (q *Query) Execute() []*Entity {
	var results []*Entity

	for _, entity := range q.world.entities {
		matches := true
		for _, id := range q.componentTypeIDs {
			if !entity.HasComponent(id) {
				matches = false

				break
			}
		}

		if matches {
			results = append(results, entity)
		}
	}

	return results
}
