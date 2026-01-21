package ecs

import "reflect"

type SystemCtor[S System] func(priority int) S

var (
	systemsRegistry map[string]SystemCtor[System] = make(map[string]SystemCtor[System])
)

func getName[S System]() string {
	return reflect.TypeFor[S]().Name()
}

func RegisterSystem[S System](systemCtor SystemCtor[S]) {
	name := getName[S]()
	systemsRegistry[name] = func(priority int) System {
		return systemCtor(priority)
	}
}

func GetSystem(name string) (SystemCtor[System], bool) {
	ctor, ok := systemsRegistry[name]
	if !ok {
		return nil, false
	}

	return ctor, true
}
