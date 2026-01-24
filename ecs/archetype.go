package ecs

import (
	"fmt"
	"iter"
	"reflect"
)

type Component interface {
	Reset()
	Init()
}

// Archetype represents a group of entities with the same component signature.
type Archetype struct {
	signature     []reflect.Type
	signatureMask uint64 // Bitmask for fast signature comparison
	entities      []EntityID
	components    map[reflect.Type][]any // Component type -> slice of component data
	entityLookup  map[EntityID]int       // Entity ID -> index in entities array
}

func NewArchetype(signature []reflect.Type, signatureMask uint64) *Archetype {
	arch := &Archetype{
		signature:     signature,
		signatureMask: signatureMask,
		entities:      make([]EntityID, 0, 64),
		components:    make(map[reflect.Type][]any),
		entityLookup:  make(map[EntityID]int),
	}

	for _, compType := range signature {
		arch.components[compType] = make([]any, 0, 64)
	}

	return arch
}

// AddEntity adds an entity with its component data to the archetype.
func (a *Archetype) AddEntity(entityID EntityID, componentsData map[reflect.Type]any) {
	if _, exists := a.entityLookup[entityID]; exists {
		return
	}

	index := len(a.entities)
	a.entities = append(a.entities, entityID)
	a.entityLookup[entityID] = index

	for _, componentType := range a.signature {
		componentData, exists := componentsData[componentType]
		if !exists {
			panic(fmt.Sprintf("Component of type %s not provided for entity %d", componentType.Name(), entityID))
		}

		a.components[componentType] = append(a.components[componentType], componentData)
	}
}

// RemoveEntity removes an entity and its component data from the archetype.
func (a *Archetype) RemoveEntity(entityID EntityID) map[reflect.Type]any {
	index, exists := a.entityLookup[entityID]
	if !exists {
		return nil
	}

	// Extract component data before removal
	componentData := make(map[reflect.Type]any)
	for _, componentType := range a.signature {
		component := a.components[componentType][index]
		componentData[componentType] = component
	}

	// Swap-and-pop removal
	lastIndex := len(a.entities) - 1
	if index != lastIndex {
		lastEntityID := a.entities[lastIndex]
		a.entities[index] = lastEntityID
		a.entityLookup[lastEntityID] = index

		for _, componentType := range a.signature {
			a.components[componentType][index] = a.components[componentType][lastIndex]
		}
	}

	a.entities = a.entities[:lastIndex]
	for _, componentType := range a.signature {
		a.components[componentType] = a.components[componentType][:lastIndex]
	}

	delete(a.entityLookup, entityID)

	return componentData
}

func (a *Archetype) GetComponent(entityID EntityID, componentType reflect.Type) (any, bool) {
	index, exists := a.entityLookup[entityID]
	if !exists {
		return nil, false
	}

	components, exists := a.components[componentType]
	if !exists {
		return nil, false
	}

	return components[index], true
}

func (a *Archetype) HasComponent(componentTpe reflect.Type) bool {
	_, exists := a.components[componentTpe]
	return exists
}

// MatchesQuery checks if the archetype matches a query based on component signatures.
func (a *Archetype) MatchesQuery(queryMask uint64) bool {
	// Subset match: Archetype must have at least all components in queryMask
	return (a.signatureMask & queryMask) == queryMask
}

func (a *Archetype) Entities() iter.Seq[EntityID] {
	return func(yield func(EntityID) bool) {
		for _, entityID := range a.entities {
			if !yield(entityID) {
				break
			}
		}
	}
}

func (a *Archetype) Count() int {
	return len(a.entities)
}

// SignatureMatches checks if two signatures are identical
func (a *Archetype) SignatureMatches(queryMask uint64) bool {
	return a.signatureMask == queryMask
}
