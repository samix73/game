package ecs

type World interface {
	Update() error
	Draw()
}
