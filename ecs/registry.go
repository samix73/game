package ecs

import (
	"log/slog"
	"reflect"
)

type (
	SystemCtor[S System] func(priority int) S
)

const maxComponents = 64

var (
	systemsRegistry    = make(map[string]SystemCtor[System])
	componentsRegistry = make(map[string]any)
	componentTypeBits  = make(map[reflect.Type]uint)
	nextComponentBit   uint
)

func getName[S any]() string {
	t := reflect.TypeFor[S]()
	if t.Kind() == reflect.Pointer {
		return t.Elem().Name()
	}

	return t.Name()
}

// RegisterSystem registers a system constructor in the ECS registry
// to allow for dynamic system creation.
func RegisterSystem[S System](systemCtor SystemCtor[S]) {
	name := getName[S]()
	if _, ok := systemsRegistry[name]; ok {
		slog.Error("ecs.RegisterSystem: system already registered", slog.String("name", name))
		return
	}

	systemsRegistry[name] = func(priority int) System {
		return systemCtor(priority)
	}

	slog.Debug("ecs.RegisterSystem: registered system", slog.String("name", name))
}

func RegisterComponent[T any]() {
	name := getName[T]()
	if _, ok := componentsRegistry[name]; ok {
		slog.Error("ecs.RegisterComponent: component already registered", slog.String("name", name))
		return
	}

	componentsRegistry[name] = *new(T)

	// Assign a unique bit position for bitmask filtering
	componentType := reflect.TypeFor[T]()
	if _, exists := componentTypeBits[componentType]; !exists {
		if nextComponentBit < maxComponents {
			componentTypeBits[componentType] = nextComponentBit
			nextComponentBit++
		} else {
			slog.Warn("ecs.RegisterComponent: exceeded 64 component types, bitmask optimization disabled for new components",
				slog.String("name", name))
		}
	}

	slog.Debug("ecs.RegisterComponent: registered component", slog.String("name", name))
}

func NewComponent(em *EntityManager, name string) (any, bool) {
	comp, ok := componentsRegistry[name]
	if !ok {
		return nil, false
	}

	componentType := reflect.TypeOf(comp)
	pool := em.getOrCreatePool(componentType, func() any {
		return reflect.New(componentType).Interface()
	})

	return pool.Get(), true
}

// GetSystem retrieves a system constructor from the ECS registry by name.
func GetSystem(name string) (SystemCtor[System], bool) {
	ctor, ok := systemsRegistry[name]
	if !ok {
		return nil, false
	}

	return ctor, true
}

// getComponentsBitmask computes a bitmask from component types for fast signature matching.
func getComponentsBitmask(componentTypes []archetypeComponentSignature) (Bitmask, bool) {
	var mask Bitmask
	for _, ct := range componentTypes {
		bitPos := ct.bit
		if bitPos == 0 {
			var ok bool
			bitPos, ok = getComponentBit(ct.typ)
			if !ok {
				return 0, false
			}
		}

		mask.SetFlag(bitPos)
	}

	return mask, true
}

func getComponentBit(componentType reflect.Type) (uint, bool) {
	bitPos, exists := componentTypeBits[componentType]
	return bitPos, exists
}

func ComponentsBitMask(head reflect.Type, tail ...reflect.Type) (Bitmask, bool) {
	componentSignature := make([]archetypeComponentSignature, 1+len(tail))
	componentSignature[0].typ = head
	for i := range tail {
		componentSignature[i+1].typ = tail[i]
	}
	return getComponentsBitmask(componentSignature)
}

func ComponentBit[T any]() (uint, bool) {
	return getComponentBit(reflect.TypeFor[T]())
}
