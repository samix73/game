package netcode

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

// NetworkVariableInterface allows the NetworkSystem to iterate generic variables without knowing the type T.
type NetworkVariableInterface interface {
	IsDirty() bool
	ResetDirty()
	Serialize() ([]byte, error)
	Deserialize(data []byte) error
}

// NetworkVariable synchronizes a value of type T over the network.
type NetworkVariable[T any] struct {
	Value T
	dirty bool
}

// Set updates the value and marks it for replication.
func (v *NetworkVariable[T]) Set(val T) {
	v.Value = val
	v.dirty = true
}

// Get returns the current value.
func (v *NetworkVariable[T]) Get() T {
	return v.Value
}

func (v *NetworkVariable[T]) IsDirty() bool {
	return v.dirty
}

func (v *NetworkVariable[T]) ResetDirty() {
	v.dirty = false
}

func (v *NetworkVariable[T]) Serialize() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(v.Value); err != nil {
		return nil, fmt.Errorf("failed to encode value: %w", err)
	}

	return buf.Bytes(), nil
}

func (v *NetworkVariable[T]) Deserialize(data []byte) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	if err := dec.Decode(&v.Value); err != nil {
		return fmt.Errorf("failed to decode value: %w", err)
	}

	return nil
}
