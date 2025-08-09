package ecs

import (
	"context"
	"iter"
	"runtime/trace"
	"sync"
)

type Component interface {
	Reset()
}

type ComponentContainer struct {
	pool sync.Pool

	components []any
	entityIDs  []EntityID

	componentLookupMap map[EntityID]int
}

func NewComponentContainer(ctx context.Context, newFn func() any) *ComponentContainer {
	ctx, task := trace.NewTask(ctx, "ecs.NewComponentContainer")
	defer task.End()

	return &ComponentContainer{
		pool: sync.Pool{New: func() any { return newFn() }},

		components:         make([]any, 0, 1024),
		entityIDs:          make([]EntityID, 0, 1024),
		componentLookupMap: make(map[EntityID]int),
	}
}

func (c *ComponentContainer) Add(ctx context.Context, entityID EntityID) any {
	ctx, task := trace.NewTask(ctx, "ecs.ComponentContainer.Add")
	defer task.End()

	if _, ok := c.componentLookupMap[entityID]; ok {
		return nil
	}

	component := c.pool.Get()

	c.components = append(c.components, component)
	c.entityIDs = append(c.entityIDs, entityID)
	c.componentLookupMap[entityID] = len(c.components) - 1

	return component
}

func (c *ComponentContainer) Remove(ctx context.Context, entityID EntityID) {
	ctx, task := trace.NewTask(ctx, "ecs.ComponentContainer.Remove")
	defer task.End()

	indexToRemove, ok := c.componentLookupMap[entityID]
	if !ok {
		return
	}

	componentToRemove := c.components[indexToRemove]

	lastIndex := len(c.entityIDs) - 1
	if lastIndex != indexToRemove {
		c.components[indexToRemove] = c.components[lastIndex]
		c.entityIDs[indexToRemove] = c.entityIDs[lastIndex]

		c.componentLookupMap[c.entityIDs[indexToRemove]] = indexToRemove
	}

	c.components = c.components[:lastIndex]
	c.entityIDs = c.entityIDs[:lastIndex]

	delete(c.componentLookupMap, entityID)

	if typedComponent, ok := componentToRemove.(Component); ok {
		typedComponent.Reset()
	}
	c.pool.Put(componentToRemove)
}

func (c *ComponentContainer) All(ctx context.Context) iter.Seq2[EntityID, any] {
	ctx, task := trace.NewTask(ctx, "ecs.ComponentContainer.All")
	defer task.End()

	return func(yield func(EntityID, any) bool) {
		for i, entityID := range c.entityIDs {
			if !yield(entityID, c.components[i]) {
				break
			}
		}
	}
}

func (c *ComponentContainer) Get(ctx context.Context, entityID EntityID) (any, bool) {
	ctx, task := trace.NewTask(ctx, "ecs.ComponentContainer.Get")
	defer task.End()

	index, ok := c.componentLookupMap[entityID]
	if !ok {
		return nil, false
	}

	return c.components[index], true
}

func (c *ComponentContainer) Count(ctx context.Context) int {
	ctx, task := trace.NewTask(ctx, "ecs.ComponentContainer.Count")
	defer task.End()

	return len(c.components)
}

func (c *ComponentContainer) Entities(ctx context.Context) iter.Seq[EntityID] {
	ctx, task := trace.NewTask(ctx, "ecs.ComponentContainer.Entities")
	defer task.End()

	return func(yield func(EntityID) bool) {
		for _, entityID := range c.entityIDs {
			if !yield(entityID) {
				break
			}
		}
	}
}

func (c *ComponentContainer) Components(ctx context.Context) iter.Seq[any] {
	ctx, task := trace.NewTask(ctx, "ecs.ComponentContainer.Components")
	defer task.End()

	return func(yield func(any) bool) {
		for _, component := range c.components {
			if !yield(component) {
				break
			}
		}
	}
}
