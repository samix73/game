package ecs

import (
	"log/slog"
	"reflect"
)

type (
	SystemCtor[S System] func(priority int) S
)

var (
	systemsRegistry    = make(map[string]SystemCtor[System])
	componentsRegistry = make(map[string]any)
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

	slog.Debug("ecs.RegisterComponent: registered component", slog.String("name", name))
}

func NewComponent(name string) (any, bool) {
	comp, ok := componentsRegistry[name]
	if !ok {
		return nil, false
	}

	v := reflect.New(reflect.TypeOf(comp))

	return v.Interface(), true
}

// GetSystem retrieves a system constructor from the ECS registry by name.
func GetSystem(name string) (SystemCtor[System], bool) {
	ctor, ok := systemsRegistry[name]
	if !ok {
		return nil, false
	}

	return ctor, true
}
