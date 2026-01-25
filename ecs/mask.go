package ecs

type Bitmask uint64

func NewBitmask(bits ...uint) Bitmask {
	var mask Bitmask
	for _, bit := range bits {
		mask.SetFlag(bit)
	}
	return mask
}

func (m *Bitmask) HasFlags(queryMask Bitmask) bool {
	return (*m & queryMask) == queryMask
}

func (m *Bitmask) HasFlag(bitPos uint) bool {
	return (*m & (1 << bitPos)) != 0
}

func (m *Bitmask) SetFlag(bitPos uint) {
	*m |= (1 << bitPos)
}

func (m *Bitmask) ClearFlag(bitPos uint) {
	*m &= ^(1 << bitPos)
}

func (m *Bitmask) Clear() {
	*m = 0
}

func (m *Bitmask) Clone() Bitmask {
	return *m
}
