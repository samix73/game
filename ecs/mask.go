package ecs

type Bitmask uint64

func (m *Bitmask) HasFlags(queryMask Bitmask) bool {
	return (*m & queryMask) == queryMask
}

func (m *Bitmask) HasFlag(bitPos uint) bool {
	return (*m & (1 << bitPos)) != 0
}

func (m *Bitmask) SetFlag(bitPos uint) {
	*m |= (1 << bitPos)
}
