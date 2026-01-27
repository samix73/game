package ecs

import (
	"bytes"
	"fmt"
	"iter"
	"log/slog"
	"reflect"
	"slices"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/samix73/game/game/assets"
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

func (em *EntityManager) NewEntity(components ...any) (EntityID, error) {
	id := em.nextID
	em.nextID++

	signature := make([]archetypeComponentSignature, 0, len(components))
	componentData := make(map[reflect.Type]any, len(components))
	for _, component := range components {
		typ := reflect.TypeOf(component)
		if typ.Kind() == reflect.Pointer {
			typ = typ.Elem()
		}

		bit, ok := getComponentBit(typ)
		if !ok {
			return 0, fmt.Errorf("ecs.EntityManager.NewEntity: component %s not registered, call RegisterComponent first", typ.Name())
		}

		signature = append(signature, archetypeComponentSignature{
			typ: typ,
			bit: bit,
		})
		componentData[typ] = component
	}

	archetype, err := em.getOrCreateArchetype(signature)
	if err != nil {
		return 0, fmt.Errorf("ecs.EntityManager.NewEntity: %w", err)
	}
	if err := archetype.AddEntity(id, componentData); err != nil {
		return 0, fmt.Errorf("ecs.EntityManager.NewEntity: %w", err)
	}
	em.entityArchetype[id] = archetype

	return id, nil
}

func (em *EntityManager) getOrCreateArchetype(componentTypes []archetypeComponentSignature) (*Archetype, error) {
	signatureMask, ok := getComponentsBitmask(componentTypes)
	if !ok {
		return nil, fmt.Errorf("ecs.EntityManager.getOrCreateArchetype: component type not registered, call RegisterComponent first")
	}

	for _, arch := range em.archetypes {
		if arch.SignatureMatches(signatureMask) {
			return arch, nil
		}
	}

	archetype, err := NewArchetype(componentTypes, signatureMask)
	if err != nil {
		return nil, fmt.Errorf("ecs.EntityManager.getOrCreateArchetype: %w", err)
	}
	em.archetypes = append(em.archetypes, archetype)

	return archetype, nil
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

func (em *EntityManager) Remove(entityID EntityID) error {
	archetype, exists := em.entityArchetype[entityID]
	if !exists {
		return fmt.Errorf("ecs.EntityManager.Remove: entity %d does not exist", entityID)
	}

	componentData, err := archetype.RemoveEntity(entityID)
	if err != nil {
		return fmt.Errorf("ecs.EntityManager.Remove: %w", err)
	}

	// Return components to pools
	for componentType, component := range componentData {
		if pool, exists := em.componentPools[componentType]; exists {
			pool.Put(component)
		}
	}

	delete(em.entityArchetype, entityID)

	return nil
}

func (em *EntityManager) RemoveComponent(entityID EntityID, componentType any) error {
	archetype, exists := em.entityArchetype[entityID]
	if !exists {
		return fmt.Errorf("ecs.EntityManager.RemoveComponent: entity %d does not exist", entityID)
	}

	removedRefType := reflect.TypeOf(componentType)
	if !archetype.HasComponent(removedRefType) {
		return fmt.Errorf("ecs.EntityManager.RemoveComponent: entity %d does not have component %s", entityID, removedRefType.Name())
	}

	componentData, err := archetype.RemoveEntity(entityID)
	if err != nil {
		return fmt.Errorf("ecs.EntityManager.RemoveComponent: %w", err)
	}

	// Get the component to return to pool
	removedComponent := componentData[removedRefType]

	// Remove the specified component type
	delete(componentData, removedRefType)

	// Return removed component to pool
	if resettable, ok := removedComponent.(Component); ok {
		resettable.Reset()
	}

	pool := em.getOrCreatePool(removedRefType, func() any { return reflect.New(removedRefType).Interface() })
	pool.Put(removedComponent)

	// Calculate new signature
	newSignature := make([]archetypeComponentSignature, 0, len(archetype.Signature())-1)
	for _, t := range archetype.Signature() {
		if t.typ == removedRefType {
			continue
		}

		newSignature = append(newSignature, t)
	}

	// Move entity to new archetype
	newArchetype, err := em.getOrCreateArchetype(newSignature)
	if err != nil {
		return fmt.Errorf("ecs.EntityManager.RemoveComponent: %w", err)
	}
	if err := newArchetype.AddEntity(entityID, componentData); err != nil {
		return fmt.Errorf("ecs.EntityManager.RemoveComponent: %w", err)
	}
	em.entityArchetype[entityID] = newArchetype

	return nil
}

func (em *EntityManager) Query(queryMask Bitmask) []EntityID {
	entities := make([]EntityID, 0)

	for i := range em.archetypes {
		if em.archetypes[i].MatchesQuery(queryMask) {
			entities = append(entities, em.archetypes[i].Entities()...)
		}
	}

	return entities
}

// LoadEntity loads an entity asset and adds it to the entity manager.
func (em *EntityManager) LoadEntity(name string) (EntityID, error) {
	entityData, err := assets.GetEntity(name)
	if err != nil {
		return 0, fmt.Errorf("ecs.EntityManager.LoadEntity: %w", err)
	}

	var protoEntity EntityComponentsConfig
	md, err := toml.NewDecoder(bytes.NewReader(entityData)).Decode(&protoEntity)
	if err != nil {
		return 0, fmt.Errorf("ecs.EntityManager.LoadEntity: %w", err)
	}

	entity, err := em.NewEntity()
	if err != nil {
		return 0, fmt.Errorf("ecs.EntityManager.LoadEntity: %w", err)
	}

	for componentName, args := range protoEntity {
		component, ok := NewComponent(em, componentName)
		if !ok {
			return 0, fmt.Errorf("ecs.EntityManager.LoadEntity: component %s not found", componentName)
		}

		if err := md.PrimitiveDecode(args, component); err != nil {
			return 0, fmt.Errorf("ecs.EntityManager.LoadEntity: PrimitiveDecode %w", err)
		}

		if err := em.AddComponent(entity, component); err != nil {
			return 0, fmt.Errorf("ecs.EntityManager.LoadEntity: AddComponent %w", err)
		}
	}

	return entity, nil
}

func (em *EntityManager) AddComponent(entityID EntityID, component any) error {
	archetype, exists := em.entityArchetype[entityID]
	if !exists {
		return fmt.Errorf("entity %d does not exist", entityID)
	}

	componentType := reflect.TypeOf(component)
	if componentType.Kind() == reflect.Pointer {
		componentType = componentType.Elem()
	}

	if archetype.HasComponent(componentType) {
		return fmt.Errorf("entity %d already has component of type %s", entityID, componentType.Name())
	}

	// Get current component data
	componentData, err := archetype.RemoveEntity(entityID)
	if err != nil {
		return fmt.Errorf("ecs.EntityManager.AddComponent: %w", err)
	}

	if resettable, ok := any(component).(Component); ok {
		resettable.Init()
	}

	// Add new component
	componentData[componentType] = component

	bitPos, ok := getComponentBit(componentType)
	if !ok {
		return fmt.Errorf("component type %s not registered, call RegisterComponent first", componentType.Name())
	}

	// Calculate new signature
	newSignature := make([]archetypeComponentSignature, 0, len(archetype.signature)+1)
	newSignature = append(newSignature, archetype.signature...)
	newSignature = append(newSignature, archetypeComponentSignature{
		typ: componentType,
		bit: bitPos,
	})

	// Move entity to new archetype
	newArchetype, err := em.getOrCreateArchetype(newSignature)
	if err != nil {
		return fmt.Errorf("ecs.EntityManager.AddComponent: %w", err)
	}
	if err := newArchetype.AddEntity(entityID, componentData); err != nil {
		return fmt.Errorf("ecs.EntityManager.AddComponent: %w", err)
	}
	em.entityArchetype[entityID] = newArchetype

	return nil
}

func (em *EntityManager) Teardown() {
	em.archetypes = nil
	em.entityArchetype = nil
	em.componentPools = nil
}

func AddComponent[C any](em *EntityManager, entityID EntityID) (*C, error) {
	componentType := reflect.TypeFor[C]()
	pool := em.getOrCreatePool(componentType, func() any {
		return new(C)
	})

	component := pool.Get().(*C)
	if err := em.AddComponent(entityID, component); err != nil {
		slog.Error("ecs.AddComponent: failed to add component",
			slog.String("type", componentType.Name()),
			slog.Uint64("entityID", entityID),
			slog.Any("error", err),
		)

		return nil, fmt.Errorf("ecs.AddComponent: failed to add component: %w", err)
	}

	return component, nil
}

func RemoveComponent[C any](em *EntityManager, entityID EntityID) error {
	var zero C
	if err := em.RemoveComponent(entityID, zero); err != nil {
		return fmt.Errorf("ecs.RemoveComponent: %w", err)
	}

	return nil
}

func Query[C any](em *EntityManager) []EntityID {
	queryMask, ok := ComponentsBitMask(reflect.TypeFor[C]())
	if !ok {
		return []EntityID{}
	}

	return em.Query(queryMask)
}

func Query2[C1, C2 any](em *EntityManager) []EntityID {
	queryMask, ok := ComponentsBitMask(
		reflect.TypeFor[C1](),
		reflect.TypeFor[C2](),
	)
	if !ok {
		return []EntityID{}
	}

	return em.Query(queryMask)
}

func Query3[C1, C2, C3 any](em *EntityManager) []EntityID {
	queryMask, ok := ComponentsBitMask(
		reflect.TypeFor[C1](),
		reflect.TypeFor[C2](),
		reflect.TypeFor[C3](),
	)
	if !ok {
		return []EntityID{}
	}

	return em.Query(queryMask)
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

	componentType := reflect.TypeFor[C]()

	dataPtr, exists := archetype.GetComponentPtr(entityID, componentType)
	if !exists {
		return nil, false
	}

	return (*C)(dataPtr), true
}

func MustGetComponent[C any](em *EntityManager, entityID EntityID) *C {
	component, exists := GetComponent[C](em, entityID)
	if !exists {
		panic(fmt.Sprintf("Entity %d does not have component of type %s", entityID, reflect.TypeFor[C]().Name()))
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

func QueryWith[C any](em *EntityManager, filter Filter[C]) []EntityID {
	if filter == nil {
		return Query[C](em)
	}

	entities := Query[C](em)
	for i, entityID := range entities {
		if !evaluateFilter(em, entityID, filter) {
			entities = append(entities[:i], entities[i+1:]...)
		}
	}

	return entities
}

func QueryWith2[C1, C2 any](em *EntityManager, filter1 Filter[C1], filter2 Filter[C2]) []EntityID {
	if filter1 == nil && filter2 == nil {
		return Query2[C1, C2](em)
	}

	entities := Query2[C1, C2](em)
	for i, entityID := range entities {
		if !evaluateFilter(em, entityID, filter1) || !evaluateFilter(em, entityID, filter2) {
			entities = append(entities[:i], entities[i+1:]...)
		}
	}

	return entities
}

func QueryWith3[C1, C2, C3 any](em *EntityManager, filter1 Filter[C1], filter2 Filter[C2], filter3 Filter[C3]) []EntityID {
	if filter1 == nil && filter2 == nil && filter3 == nil {
		return Query3[C1, C2, C3](em)
	}

	entities := Query3[C1, C2, C3](em)
	for i, entityID := range entities {
		if !evaluateFilter(em, entityID, filter1) || !evaluateFilter(em, entityID, filter2) || !evaluateFilter(em, entityID, filter3) {
			entities = append(entities[:i], entities[i+1:]...)
		}
	}

	return entities
}
