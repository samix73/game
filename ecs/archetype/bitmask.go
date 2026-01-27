package ecs

type Bitmask []uint64

func NewBitmask() Bitmask {
	return Bitmask{}
}

// Set sets the i-th bit
func (b *Bitmask) Set(i uint) {
	word, bit := i/64, i%64

	if word >= uint(len(*b)) {
		*b = append(*b, make([]uint64, word-uint(len(*b))+1)...)
	}

	(*b)[word] |= 1 << bit
}

// Unset unsets the i-th bit
func (b *Bitmask) Unset(i uint) {
	word, bit := i/64, i%64

	if word < uint(len(*b)) {
		(*b)[word] &^= 1 << bit
	}
}

// HasAll returns true if all bits are set
func (b *Bitmask) HasAll(bits Bitmask) bool {
	for i, word := range bits {
		if word == 0 {
			continue
		}
		if i >= len(*b) {
			return false
		}
		if ((*b)[i] & word) != word {
			return false
		}
	}
	return true
}

// Has returns true if the i-th bit is set
func (b *Bitmask) Has(i uint) bool {
	word, bit := i/64, i%64
	return word < uint(len(*b)) && (*b)[word]&(1<<bit) != 0
}
