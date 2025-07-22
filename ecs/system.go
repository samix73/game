package ecs

type System interface {
	Priority() int
	Update() error
	Draw()
	Remove()
}
