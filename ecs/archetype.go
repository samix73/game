package ecs

import (
	"fmt"
	"iter"
	"reflect"
	"unsafe"
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
	signatureMask Bitmask // Bitmask for fast signature comparison
	signature     []archetypeComponentSignature
	entities      []EntityID
	components    map[uint][]byte  // Indexed by component bit position
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
		components:    make(map[uint][]byte), // Support up to 64 components
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
		bitPos := componentType.bit

		componentData, exists := componentsData[typ]
		if !exists {
			return fmt.Errorf("Component of type %s not provided for entity %d", typ.Name(), entityID)
		}

		if componentData == nil {
			return fmt.Errorf("Component of type %s is nil for entity %d", typ.Name(), entityID)
		}

		val := reflect.ValueOf(componentData)
		if val.Kind() != reflect.Pointer {
			return fmt.Errorf("Component of type %s must be a pointer", typ.Name())
		}

		dataPtr := val.UnsafePointer()

		src := unsafe.Slice((*byte)(dataPtr), typ.Size())
		a.components[bitPos] = append(a.components[bitPos], src...)
	}

	return nil
}

// RemoveEntity removes an entity and its component data from the archetype.
func (a *Archetype) RemoveEntity(entityID EntityID) (map[reflect.Type]any, error) {
	index, exists := a.entityLookup[entityID]
	if !exists {
		return nil, fmt.Errorf("entity %d not found in archetype", entityID)
	}

	// Extract component data before removal
	componentData := make(map[reflect.Type]any)
	for _, componentType := range a.signature {
		typ := componentType.typ
		bitPos := componentType.bit

		// Calculate the offset for this entity's component data
		offset := uintptr(index) * typ.Size()
		buffer, exists := a.components[bitPos]

		// Check if buffer exists and has sufficient length
		if !exists {
			return nil, fmt.Errorf("ecs.Archetype.RemoveEntity: Component of type %s not found for entity %d", typ.Name(), entityID)
		}

		if len(buffer) == 0 {
			return nil, fmt.Errorf("ecs.Archetype.RemoveEntity: Component of type %s has no data for entity %d", typ.Name(), entityID)
		}

		// Create a pointer to the component data in the buffer
		dataPtr := unsafe.Pointer(&buffer[offset])

		componentData[typ] = reflect.NewAt(typ, dataPtr).Interface()
	}

	// Swap-and-pop removal
	lastIndex := len(a.entities) - 1
	if index != lastIndex {
		lastEntityID := a.entities[lastIndex]
		a.entities[index] = lastEntityID
		a.entityLookup[lastEntityID] = index

		// Copy the last entity's component data to the removed entity's position
		for _, componentType := range a.signature {
			typ := componentType.typ
			bitPos := componentType.bit
			componentSize := typ.Size()

			srcOffset := uintptr(lastIndex) * componentSize
			dstOffset := uintptr(index) * componentSize

			buffer := a.components[bitPos]
			copy(buffer[dstOffset:dstOffset+componentSize], buffer[srcOffset:srcOffset+componentSize])
		}
	}

	a.entities = a.entities[:lastIndex]
	for _, componentType := range a.signature {
		typ := componentType.typ
		componentSize := typ.Size()
		newLen := uintptr(lastIndex) * componentSize
		a.components[componentType.bit] = a.components[componentType.bit][:newLen]
	}

	delete(a.entityLookup, entityID)

	return componentData, nil
}

// GetComponentPtr returns a pointer to the component data for the given entity.
// It is the responsibility of the caller to ensure that the component type is correct
// and that the entity exists in the archetype.
func (a *Archetype) GetComponentPtr(entityID EntityID, componentType reflect.Type) (unsafe.Pointer, bool) {
	if componentType.Kind() != reflect.Struct {
		return nil, false
	}

	bitPos, exists := getComponentBit(componentType)
	if !exists {
		return nil, false
	}

	index, exists := a.entityLookup[entityID]
	if !exists {
		return nil, false
	}

	buffer, ok := a.components[bitPos]
	if !ok {
		return nil, false
	}

	offset := uintptr(index) * componentType.Size()
	dataPtr := unsafe.Pointer(&buffer[offset])

	return dataPtr, true
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
