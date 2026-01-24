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

type archetypeComponentSignature struct {
	typ reflect.Type
	bit uint // Pre-calculated bit positions for signature
}

// Archetype represents a group of entities with the same component signature.
type Archetype struct {
	signature     []archetypeComponentSignature
	signatureMask Bitmask // Bitmask for fast signature comparison
	entities      []EntityID
	components    map[uint][]any   // Indexed by component bit position
	entityLookup  map[EntityID]int // Entity ID -> index in entities array
}

func NewArchetype(componentTypes []archetypeComponentSignature, signatureMask Bitmask) (*Archetype, error) {
	signature := make([]archetypeComponentSignature, len(componentTypes))
	for i, compType := range componentTypes {
		if compType.bit == 0 {
			var exists bool
			compType.bit, exists = getComponentBit(compType.typ)
			if !exists {
				return nil, fmt.Errorf("component type %s not registered, call RegisterComponent first", compType.typ.Name())
			}
		}

		signature[i] = compType
	}

	return &Archetype{
		signature:     signature,
		signatureMask: signatureMask,
		entities:      make([]EntityID, 0, 64),
		components:    make(map[uint][]any), // Support up to 64 components
		entityLookup:  make(map[EntityID]int),
	}, nil
}

func (a *Archetype) Signature() []archetypeComponentSignature {
	return a.signature
}

// AddEntity adds an entity with its component data to the archetype.
func (a *Archetype) AddEntity(entityID EntityID, componentsData map[reflect.Type]any) error {
	if _, exists := a.entityLookup[entityID]; exists {
		return fmt.Errorf("entity %d already exists in archetype", entityID)
	}

	index := len(a.entities)
	a.entities = append(a.entities, entityID)
	a.entityLookup[entityID] = index

	for _, componentType := range a.signature {
		typ := componentType.typ
		componentData, exists := componentsData[typ]
		if !exists {
			return fmt.Errorf("Component of type %s not provided for entity %d", typ.Name(), entityID)
		}

		bitPos := componentType.bit
		a.components[bitPos] = append(a.components[bitPos], componentData)
	}

	return nil
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
		bitPos := componentType.bit
		component := a.components[bitPos][index]
		componentData[componentType.typ] = component
	}

	// Swap-and-pop removal
	lastIndex := len(a.entities) - 1
	if index != lastIndex {
		lastEntityID := a.entities[lastIndex]
		a.entities[index] = lastEntityID
		a.entityLookup[lastEntityID] = index

		for _, componentType := range a.signature {
			a.components[componentType.bit][index] = a.components[componentType.bit][lastIndex]
		}
	}

	a.entities = a.entities[:lastIndex]
	for _, componentType := range a.signature {
		a.components[componentType.bit] = a.components[componentType.bit][:lastIndex]
	}

	delete(a.entityLookup, entityID)

	return componentData
}

func (a *Archetype) GetComponent(entityID EntityID, componentType reflect.Type) (any, bool) {
	bitPos, exists := getComponentBit(componentType)
	if !exists {
		return nil, false
	}

	return a.GetComponentByBit(entityID, bitPos)
}

func (a *Archetype) GetComponentByBit(entityID EntityID, bitPos uint) (any, bool) {
	index, exists := a.entityLookup[entityID]
	if !exists {
		return nil, false
	}

	if !a.signatureMask.HasFlag(bitPos) {
		return nil, false
	}

	return a.components[bitPos][index], true
}

func (a *Archetype) HasComponent(componentType reflect.Type) bool {
	bitPos, exists := getComponentBit(componentType)
	if !exists {
		return false
	}

	return a.signatureMask.HasFlag(bitPos)
}

// MatchesQuery checks if the archetype matches a query based on component signatures.
func (a *Archetype) MatchesQuery(queryMask Bitmask) bool {
	// Subset match: Archetype must have at least all components in queryMask
	return a.signatureMask.HasFlags(queryMask)
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
func (a *Archetype) SignatureMatches(queryMask Bitmask) bool {
	return a.signatureMask == queryMask
}
