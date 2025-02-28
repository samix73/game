package data

type EntityType int

const (
	EntityTypePlayer EntityType = iota + 1
)

type EntityID uint64

type Entity interface {
	ID() EntityID
	SetID(id EntityID)
	Type() EntityType
}
