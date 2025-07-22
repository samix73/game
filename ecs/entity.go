package ecs

import (
	"iter"
	"reflect"
)

type EntityID = ID

type EntityManager struct {
	entities                  map[EntityID]struct{}
	componentContainers       map[reflect.Type]*ComponentContainer
	entityComponentSignatures map[EntityID]map[reflect.Type]struct{}
}

func NewEntityManager() *EntityManager {
	return &EntityManager{
		entities:                  make(map[EntityID]struct{}),
		componentContainers:       make(map[reflect.Type]*ComponentContainer),
		entityComponentSignatures: make(map[EntityID]map[reflect.Type]struct{}),
	}
}

func (em *EntityManager) NewEntity() EntityID {
	id := NextID()
	em.entities[id] = struct{}{}
	em.entityComponentSignatures[id] = make(map[reflect.Type]struct{})

	return id
}

func (em *EntityManager) AddComponent(entityID EntityID, component any) {
	// Check if the entity exists
	if _, exists := em.entities[entityID]; !exists {
		return
	}

	// Check if the component type is already registered for this entity
	componentType := reflect.TypeOf(component)
	if _, exists := em.entityComponentSignatures[entityID][componentType]; exists {
		return
	}

	container, exists := em.componentContainers[componentType]
	if !exists {
		container = NewComponentContainer()
		em.componentContainers[componentType] = container
	}

	container.Add(entityID, component)
	em.entityComponentSignatures[entityID][componentType] = struct{}{}
}

func (em *EntityManager) GetComponent(entityID EntityID, componentType reflect.Type) (any, bool) {
	if _, exists := em.entities[entityID]; !exists {
		return nil, false
	}

	if _, exists := em.entityComponentSignatures[entityID][componentType]; !exists {
		return nil, false
	}

	container, exists := em.componentContainers[componentType]
	if !exists {
		return nil, false
	}

	component, exists := container.Get(entityID)
	if !exists {
		return nil, false
	}

	return component, true
}

func (em *EntityManager) HasComponent(entityID EntityID, componentType reflect.Type) bool {
	if _, exists := em.entities[entityID]; !exists {
		return false
	}

	if _, exists := em.entityComponentSignatures[entityID][componentType]; !exists {
		return false
	}

	return true
}

func (em *EntityManager) Remove(entityID EntityID) {
	if _, exists := em.entities[entityID]; !exists {
		return
	}

	for componentType := range em.entityComponentSignatures[entityID] {
		if container, exists := em.componentContainers[componentType]; exists {
			container.Remove(entityID)
		}
	}

	delete(em.entityComponentSignatures, entityID)
	delete(em.entities, entityID)
}

// Query returns a sequence of EntityIDs that match the specified component types.
func (em *EntityManager) Query(componentTypes ...reflect.Type) iter.Seq[EntityID] {
	zeroIter := func(yield func(EntityID) bool) {}

	if len(componentTypes) == 0 {
		return zeroIter
	}

	// If only one component type is specified, return entities with that component
	if len(componentTypes) == 1 {
		componentContainer, exists := em.componentContainers[componentTypes[0]]
		if !exists {
			return zeroIter
		}

		return componentContainer.Entities()
	}

	// Pre-check: if any component type doesn't exist, return empty iterator
	containers := make([]*ComponentContainer, len(componentTypes))
	for i, componentType := range componentTypes {
		container, exists := em.componentContainers[componentType]
		if !exists {
			return zeroIter
		}
		containers[i] = container
	}

	// Find the container with the smallest number of entities to start with
	// This reduces the number of entities we need to check
	smallestIdx := 0
	smallestCount := containers[0].Count()

	// Check if the smallest container is empty
	if smallestCount == 0 {
		return zeroIter
	}

	for i := 1; i < len(containers); i++ {
		count := containers[i].Count()
		if count == 0 {
			// If any container has zero entities, we can return immediately
			return zeroIter
		} else if count < smallestCount {
			smallestCount = count
			smallestIdx = i
		}
	}

	// Start with the smallest set and filter iteratively
	smallestContainer := containers[smallestIdx]
	otherContainers := make([]*ComponentContainer, 0, len(containers)-1)
	for i, container := range containers {
		if i != smallestIdx {
			otherContainers = append(otherContainers, container)
		}
	}

	return func(yield func(EntityID) bool) {
		for entityID := range smallestContainer.Entities() {
			// Check if this entity exists in all other containers
			hasAllComponents := true
			for _, container := range otherContainers {
				if _, exists := container.Get(entityID); !exists {
					hasAllComponents = false
					break
				}
			}

			if hasAllComponents {
				if !yield(entityID) {
					break
				}
			}
		}
	}
}

func Query[C any](em *EntityManager) iter.Seq[EntityID] {
	return em.Query(reflect.TypeOf((*C)(nil)))
}

func Query2[C1, C2 any](em *EntityManager) iter.Seq[EntityID] {
	return em.Query(reflect.TypeOf((*C1)(nil)), reflect.TypeOf((*C2)(nil)))
}

func Query3[C1, C2, C3 any](em *EntityManager) iter.Seq[EntityID] {
	return em.Query(reflect.TypeOf((*C1)(nil)), reflect.TypeOf((*C2)(nil)), reflect.TypeOf((*C3)(nil)))
}

func HasComponent[C any](em *EntityManager, entityID EntityID) bool {
	return em.HasComponent(entityID, reflect.TypeOf((*C)(nil)))
}

func GetComponent[C any](em *EntityManager, entityID EntityID) (*C, bool) {
	component, exists := em.GetComponent(entityID, reflect.TypeOf((*C)(nil)))
	if !exists {
		return nil, false
	}
	return component.(*C), true
}
