package ecs

import (
	"fmt"
	"iter"
	"reflect"
	"slices"
	"sync"
)

type EntityID = uint64

type EntityManager struct {
	nextID          EntityID
	archetypes      []*Archetype
	entityArchetype map[EntityID]*Archetype
	componentPools  map[reflect.Type]*sync.Pool
}

func NewEntityManager() *EntityManager {
	return &EntityManager{
		nextID:          1,
		archetypes:      make([]*Archetype, 0),
		entityArchetype: make(map[EntityID]*Archetype),
		componentPools:  make(map[reflect.Type]*sync.Pool),
	}
}

func (em *EntityManager) NewEntity() EntityID {
	id := em.nextID
	em.nextID++

	archetype := em.getOrCreateArchetype([]reflect.Type{})
	archetype.AddEntity(id, make(map[reflect.Type]any))
	em.entityArchetype[id] = archetype

	return id
}

func (em *EntityManager) getOrCreateArchetype(componentTypes []reflect.Type) *Archetype {
	for _, arch := range em.archetypes {
		if arch.SignatureMatches(componentTypes) {
			return arch
		}
	}

	archetype := NewArchetype(componentTypes)
	em.archetypes = append(em.archetypes, archetype)

	return archetype
}

func (em *EntityManager) getOrCreatePool(componentType reflect.Type, newFn func() any) *sync.Pool {
	if pool, exists := em.componentPools[componentType]; exists {
		return pool
	}

	pool := &sync.Pool{New: newFn}
	em.componentPools[componentType] = pool

	return pool
}

func (em *EntityManager) HasComponent(entityID EntityID, componentType any) bool {
	archetype, exists := em.entityArchetype[entityID]
	if !exists {
		return false
	}

	return archetype.HasComponent(reflect.TypeOf(componentType))
}

func (em *EntityManager) Remove(entityID EntityID) {
	archetype, exists := em.entityArchetype[entityID]
	if !exists {
		return
	}

	componentData := archetype.RemoveEntity(entityID)

	// Return components to pools
	for componentType, component := range componentData {
		if pool, exists := em.componentPools[componentType]; exists {
			pool.Put(component)
		}
	}

	delete(em.entityArchetype, entityID)
}

func (em *EntityManager) RemoveComponent(entityID EntityID, componentType any) {
	archetype, exists := em.entityArchetype[entityID]
	if !exists {
		return
	}

	refType := reflect.TypeOf(componentType)
	if !archetype.HasComponent(refType) {
		return
	}

	componentData := archetype.RemoveEntity(entityID)

	// Get the component to return to pool
	removedComponent := componentData[refType]

	// Remove the specified component type
	delete(componentData, refType)

	// Return removed component to pool
	if pool, exists := em.componentPools[refType]; exists {
		if resettable, ok := removedComponent.(Component); ok {
			resettable.Reset()
		}
		pool.Put(removedComponent)
	}

	// Calculate new signature
	newSignature := make([]reflect.Type, 0, len(archetype.signature)-1)
	for _, t := range archetype.signature {
		if t != refType {
			newSignature = append(newSignature, t)
		}
	}

	// Move entity to new archetype
	newArchetype := em.getOrCreateArchetype(newSignature)
	newArchetype.AddEntity(entityID, componentData)
	em.entityArchetype[entityID] = newArchetype
}

func (em *EntityManager) Query(componentTypes ...any) iter.Seq[EntityID] {
	if len(componentTypes) == 0 {
		return func(yield func(EntityID) bool) {}
	}

	reflectTypes := make([]reflect.Type, len(componentTypes))
	for i, ct := range componentTypes {
		reflectTypes[i] = reflect.TypeOf(ct)
	}

	return func(yield func(EntityID) bool) {
		for _, archetype := range em.archetypes {
			if !archetype.MatchesQuery(reflectTypes) {
				continue
			}

			for entityID := range archetype.Entities() {
				if !yield(entityID) {
					return
				}
			}
		}
	}
}

func (em *EntityManager) AddComponent(entityID EntityID, component any) error {
	archetype, exists := em.entityArchetype[entityID]
	if !exists {
		return fmt.Errorf("entity %d does not exist", entityID)
	}

	typ := reflect.TypeOf(component)
	if typ.Kind() != reflect.Pointer {
		return fmt.Errorf("component must be a pointer type")
	}

	componentType := reflect.TypeOf(component).Elem()

	if resettable, ok := any(component).(Component); ok {
		resettable.Init()
	}

	// Get current component data
	componentData := archetype.RemoveEntity(entityID)

	// Add new component
	componentData[componentType] = component

	// Calculate new signature
	newSignature := make([]reflect.Type, 0, len(archetype.signature)+1)
	newSignature = append(newSignature, archetype.signature...)
	newSignature = append(newSignature, componentType)

	// Move entity to new archetype
	newArchetype := em.getOrCreateArchetype(newSignature)
	newArchetype.AddEntity(entityID, componentData)
	em.entityArchetype[entityID] = newArchetype

	return nil
}

func (em *EntityManager) Teardown() {
	em.archetypes = nil
	em.entityArchetype = nil
	em.componentPools = nil
}

func AddComponent[C any](em *EntityManager, entityID EntityID) *C {
	archetype, exists := em.entityArchetype[entityID]
	if !exists {
		return nil
	}

	var zero C
	componentType := reflect.TypeOf(zero)

	// Check if entity already has this component
	if archetype.HasComponent(componentType) {
		component, _ := archetype.GetComponent(entityID, componentType)
		return component.(*C)
	}

	// Get or create pool
	pool := em.getOrCreatePool(componentType, func() any {
		return new(C)
	})

	// Get component from pool
	component := pool.Get().(*C)
	if resettable, ok := any(component).(Component); ok {
		resettable.Init()
	}

	// Get current component data
	componentData := archetype.RemoveEntity(entityID)

	// Add new component
	componentData[componentType] = component

	// Calculate new signature
	newSignature := make([]reflect.Type, 0, len(archetype.signature)+1)
	newSignature = append(newSignature, archetype.signature...)
	newSignature = append(newSignature, componentType)

	// Move entity to new archetype
	newArchetype := em.getOrCreateArchetype(newSignature)
	newArchetype.AddEntity(entityID, componentData)
	em.entityArchetype[entityID] = newArchetype

	return component
}

func RemoveComponent[C any](em *EntityManager, entityID EntityID) {
	var zero C
	em.RemoveComponent(entityID, zero)
}

func Query[C any](em *EntityManager) iter.Seq[EntityID] {
	var zero C
	return em.Query(zero)
}

func Query2[C1, C2 any](em *EntityManager) iter.Seq[EntityID] {
	var zero1 C1
	var zero2 C2
	return em.Query(zero1, zero2)
}

func Query3[C1, C2, C3 any](em *EntityManager) iter.Seq[EntityID] {
	var zero1 C1
	var zero2 C2
	var zero3 C3
	return em.Query(zero1, zero2, zero3)
}

func HasComponent[C any](em *EntityManager, entityID EntityID) bool {
	var zero C
	return em.HasComponent(entityID, zero)
}

func GetComponent[C any](em *EntityManager, entityID EntityID) (*C, bool) {
	archetype, exists := em.entityArchetype[entityID]
	if !exists {
		return nil, false
	}

	var zero C
	componentType := reflect.TypeOf(zero)

	component, exists := archetype.GetComponent(entityID, componentType)
	if !exists {
		return nil, false
	}

	return component.(*C), true
}

func MustGetComponent[C any](em *EntityManager, entityID EntityID) *C {
	component, exists := GetComponent[C](em, entityID)
	if !exists {
		var zero C
		panic(fmt.Sprintf("Entity %d does not have component of type %s", entityID, reflect.TypeOf(zero).Name()))
	}

	return component
}

func First(iterator iter.Seq[EntityID]) (EntityID, bool) {
	for item := range iterator {
		return item, true
	}

	return 0, false
}

func Count(it iter.Seq[EntityID]) int {
	return len(slices.Collect(it))
}

func evaluateFilter[C any](em *EntityManager, entityID EntityID, filter Filter[C]) bool {
	if filter == nil {
		return true
	}

	component, ok := GetComponent[C](em, entityID)
	if !ok {
		return false
	}

	return filter(component)
}

func QueryWith[C any](em *EntityManager, filter Filter[C]) iter.Seq[EntityID] {
	if filter == nil {
		return Query[C](em)
	}

	return func(yield func(EntityID) bool) {
		for entityID := range Query[C](em) {
			if evaluateFilter(em, entityID, filter) {
				if !yield(entityID) {
					break
				}
			}
		}
	}
}

func QueryWith2[C1, C2 any](em *EntityManager, filter1 Filter[C1], filter2 Filter[C2]) iter.Seq[EntityID] {
	if filter1 == nil && filter2 == nil {
		return Query2[C1, C2](em)
	}

	return func(yield func(EntityID) bool) {
		for entityID := range Query2[C1, C2](em) {
			if evaluateFilter(em, entityID, filter1) && evaluateFilter(em, entityID, filter2) {
				if !yield(entityID) {
					break
				}
			}
		}
	}
}

func QueryWith3[C1, C2, C3 any](em *EntityManager, filter1 Filter[C1], filter2 Filter[C2], filter3 Filter[C3]) iter.Seq[EntityID] {
	if filter1 == nil && filter2 == nil && filter3 == nil {
		return Query3[C1, C2, C3](em)
	}

	return func(yield func(EntityID) bool) {
		for entityID := range Query3[C1, C2, C3](em) {
			if evaluateFilter(em, entityID, filter1) &&
				evaluateFilter(em, entityID, filter2) &&
				evaluateFilter(em, entityID, filter3) {
				if !yield(entityID) {
					break
				}
			}
		}
	}
}
