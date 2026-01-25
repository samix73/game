package ecs

import (
	"reflect"
	"slices"
	"testing"

	"github.com/samix73/game/helpers"
	"github.com/stretchr/testify/require"
)

func TestNewArchetype(t *testing.T) {
	type Position struct {
		X, Y, Z float64
	}
	err := RegisterComponent[Position]()
	require.NoError(t, err)

	positionBit, exists := getComponentBit(reflect.TypeFor[Position]())
	require.True(t, exists)

	tests := []struct {
		name           string
		componentTypes []archetypeComponentSignature
		signatureMask  Bitmask
		want           *Archetype
		wantErr        bool
	}{
		{
			name:           "valid: no components",
			componentTypes: []archetypeComponentSignature{},
			signatureMask:  0,
			want: &Archetype{
				signatureMask: 0,
				signature:     []archetypeComponentSignature{},
				entities:      []EntityID{},
				components:    map[uint][]byte{},
				entityLookup:  map[EntityID]int{},
			},
			wantErr: false,
		},
		{
			name: "valid: one component",
			componentTypes: []archetypeComponentSignature{
				{
					typ: reflect.TypeFor[int](),
					bit: 1,
				},
			},
			signatureMask: 1,
			want: &Archetype{
				signatureMask: 1,
				signature:     []archetypeComponentSignature{{typ: reflect.TypeFor[int](), bit: 1}},
				entities:      []EntityID{},
				components:    map[uint][]byte{},
				entityLookup:  map[EntityID]int{},
			},
			wantErr: false,
		},
		{
			name: "valid: two components",
			componentTypes: []archetypeComponentSignature{
				{
					typ: reflect.TypeFor[int](),
					bit: 1,
				},
				{
					typ: reflect.TypeFor[float64](),
					bit: 2,
				},
			},
			signatureMask: 3,
			want: &Archetype{
				signatureMask: 3,
				signature: []archetypeComponentSignature{
					{typ: reflect.TypeFor[int](), bit: 1},
					{typ: reflect.TypeFor[float64](), bit: 2},
				},
				entities:     []EntityID{},
				components:   map[uint][]byte{},
				entityLookup: map[EntityID]int{},
			},
			wantErr: false,
		},
		{
			name: "valid: no bit set",
			componentTypes: []archetypeComponentSignature{
				{
					typ: reflect.TypeFor[Position](),
					bit: 0, // should be set by NewArchetype
				},
			},
			signatureMask: 1,
			want: &Archetype{
				signatureMask: 1,
				signature:     []archetypeComponentSignature{{typ: reflect.TypeFor[Position](), bit: positionBit}},
				entities:      []EntityID{},
				components:    map[uint][]byte{},
				entityLookup:  map[EntityID]int{},
			},
			wantErr: false,
		},
		{
			name: "invalid: pointer component",
			componentTypes: []archetypeComponentSignature{
				{
					typ: reflect.TypeFor[*int](),
					bit: 0,
				},
			},
			signatureMask: 0,
			wantErr:       true,
		},
		{
			name: "invalid: component type not registered",
			componentTypes: []archetypeComponentSignature{
				{
					typ: reflect.TypeFor[int](),
					bit: 0,
				},
			},
			signatureMask: 0,
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := NewArchetype(tt.componentTypes, tt.signatureMask)
			if tt.wantErr {
				require.Error(t, gotErr)
				return
			}
			require.NoError(t, gotErr)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestArchetype_AddEntity(t *testing.T) {
	err := RegisterComponent[int]()
	require.NoError(t, err)

	intBit, ok := getComponentBit(reflect.TypeFor[int]())
	require.True(t, ok)

	err = RegisterComponent[float64]()
	require.NoError(t, err)

	float64Bit, ok := getComponentBit(reflect.TypeFor[float64]())
	require.True(t, ok)

	type entityToAdd struct {
		entityID   EntityID
		components map[reflect.Type]any
		wantErr    bool
	}

	tests := []struct {
		name           string
		componentTypes []archetypeComponentSignature
		signatureMask  Bitmask
		entities       []entityToAdd
	}{
		{
			name:           "valid: no components",
			componentTypes: []archetypeComponentSignature{},
			signatureMask:  0,
			entities: []entityToAdd{
				{
					entityID:   0,
					components: map[reflect.Type]any{},
					wantErr:    false,
				},
			},
		},
		{
			name: "valid: one component",
			componentTypes: []archetypeComponentSignature{
				{
					typ: reflect.TypeFor[int](),
					bit: intBit,
				},
			},
			signatureMask: Bitmask(intBit),
			entities: []entityToAdd{
				{
					entityID:   0,
					components: map[reflect.Type]any{reflect.TypeFor[int](): helpers.New(1)},
					wantErr:    false,
				},
			},
		},
		{
			name: "valid: two components",
			componentTypes: []archetypeComponentSignature{
				{
					typ: reflect.TypeFor[int](),
					bit: intBit,
				},
				{
					typ: reflect.TypeFor[float64](),
					bit: float64Bit,
				},
			},
			signatureMask: NewBitmask(intBit, float64Bit),
			entities: []entityToAdd{
				{
					entityID:   0,
					components: map[reflect.Type]any{reflect.TypeFor[int](): helpers.New(1), reflect.TypeFor[float64](): helpers.New(2.0)},
					wantErr:    false,
				},
			},
		},
		{
			name: "invalid: pointer component",
			componentTypes: []archetypeComponentSignature{
				{
					typ: reflect.TypeFor[int](),
					bit: intBit,
				},
			},
			signatureMask: NewBitmask(intBit),
			entities: []entityToAdd{
				{
					entityID:   0,
					components: map[reflect.Type]any{reflect.TypeFor[*int](): helpers.New(1)}, // componentsData key must be not be of a pointer type
					wantErr:    true,
				},
			},
		},
		{
			name: "invalid: nil component",
			componentTypes: []archetypeComponentSignature{
				{
					typ: reflect.TypeFor[int](),
					bit: intBit,
				},
			},
			signatureMask: NewBitmask(intBit),
			entities: []entityToAdd{
				{
					entityID:   0,
					components: map[reflect.Type]any{reflect.TypeFor[int](): nil},
					wantErr:    true,
				},
			},
		},
		{
			name: "invalid: component type not registered",
			componentTypes: []archetypeComponentSignature{
				{
					typ: reflect.TypeFor[int](),
					bit: intBit,
				},
			},
			signatureMask: NewBitmask(intBit),
			entities: []entityToAdd{
				{
					entityID:   0,
					components: map[reflect.Type]any{reflect.TypeFor[string](): helpers.New("test")},
					wantErr:    true,
				},
			},
		},
		{
			name: "invalid: same entity ID",
			componentTypes: []archetypeComponentSignature{
				{
					typ: reflect.TypeFor[int](),
					bit: intBit,
				},
			},
			signatureMask: NewBitmask(intBit),
			entities: []entityToAdd{
				{
					entityID:   0,
					components: map[reflect.Type]any{reflect.TypeFor[int](): helpers.New(1)},
					wantErr:    false,
				},
				{
					entityID:   0,
					components: map[reflect.Type]any{reflect.TypeFor[int](): helpers.New(2)},
					wantErr:    true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, err := NewArchetype(tt.componentTypes, tt.signatureMask)
			require.NoError(t, err)

			for _, entityToAdd := range tt.entities {
				gotErr := a.AddEntity(entityToAdd.entityID, entityToAdd.components)
				if entityToAdd.wantErr {
					require.Error(t, gotErr)
					return
				}
				require.NoError(t, gotErr)

				require.True(t, slices.ContainsFunc(a.entities, func(e EntityID) bool {
					return e == entityToAdd.entityID
				}))

				require.Contains(t, a.entityLookup, entityToAdd.entityID)

				for componentType, componentData := range entityToAdd.components {
					dataPtr, exists := a.GetComponentPtr(entityToAdd.entityID, componentType)
					require.Truef(t, exists, "component %v not found", componentType)

					require.Equal(t, componentData, reflect.NewAt(componentType, dataPtr).Interface())
				}
			}
		})
	}
}

func TestArchetype_RemoveEntity(t *testing.T) {
	err := RegisterComponent[int]()
	require.NoError(t, err)

	intBit, ok := getComponentBit(reflect.TypeFor[int]())
	require.True(t, ok)

	err = RegisterComponent[float64]()
	require.NoError(t, err)

	float64Bit, ok := getComponentBit(reflect.TypeFor[float64]())
	require.True(t, ok)

	type entityToAdd struct {
		entityID       EntityID
		componentsData map[reflect.Type]any
	}

	type entityToRemove struct {
		entityID EntityID
		want     map[reflect.Type]any
		wantErr  bool
	}

	tests := []struct {
		name             string
		componentTypes   []archetypeComponentSignature
		entitiesToAdd    []entityToAdd
		entitiesToRemove []entityToRemove
		signatureMask    Bitmask
	}{
		{
			name:           "valid: no components",
			componentTypes: []archetypeComponentSignature{},
			signatureMask:  0,
			entitiesToAdd: []entityToAdd{
				{
					entityID:       0,
					componentsData: map[reflect.Type]any{},
				},
			},
			entitiesToRemove: []entityToRemove{
				{
					entityID: 0,
					want:     map[reflect.Type]any{},
					wantErr:  false,
				},
			},
		},
		{
			name:           "valid: one component",
			componentTypes: []archetypeComponentSignature{{typ: reflect.TypeFor[int](), bit: intBit}},
			signatureMask:  NewBitmask(intBit),
			entitiesToAdd: []entityToAdd{
				{
					entityID: 0,
					componentsData: map[reflect.Type]any{
						reflect.TypeFor[int](): helpers.New(1),
					},
				},
			},
			entitiesToRemove: []entityToRemove{
				{
					entityID: 0,
					want: map[reflect.Type]any{
						reflect.TypeFor[int](): helpers.New(1),
					},
					wantErr: false,
				},
			},
		},
		{
			name: "valid: two components",
			componentTypes: []archetypeComponentSignature{
				{typ: reflect.TypeFor[int](), bit: intBit},
				{typ: reflect.TypeFor[float64](), bit: float64Bit},
			},
			signatureMask: NewBitmask(intBit, float64Bit),
			entitiesToAdd: []entityToAdd{
				{
					entityID: 0,
					componentsData: map[reflect.Type]any{
						reflect.TypeFor[int]():     helpers.New(1),
						reflect.TypeFor[float64](): helpers.New(2.0),
					},
				},
			},
			entitiesToRemove: []entityToRemove{
				{
					entityID: 0,
					want: map[reflect.Type]any{
						reflect.TypeFor[int]():     helpers.New(1),
						reflect.TypeFor[float64](): helpers.New(2.0),
					},
					wantErr: false,
				},
			},
		},
		{
			name: "invalid: entity not found",
			componentTypes: []archetypeComponentSignature{
				{typ: reflect.TypeFor[int](), bit: intBit},
			},
			signatureMask: NewBitmask(intBit),
			entitiesToAdd: []entityToAdd{
				{
					entityID: 0,
					componentsData: map[reflect.Type]any{
						reflect.TypeFor[int](): helpers.New(1),
					},
				},
			},
			entitiesToRemove: []entityToRemove{
				{
					entityID: 1,
					want:     nil,
					wantErr:  true,
				},
			},
		},
		{
			name: "invalid: remove same entity twice",
			componentTypes: []archetypeComponentSignature{
				{typ: reflect.TypeFor[int](), bit: intBit},
			},
			signatureMask: NewBitmask(intBit),
			entitiesToAdd: []entityToAdd{
				{
					entityID: 0,
					componentsData: map[reflect.Type]any{
						reflect.TypeFor[int](): helpers.New(1),
					},
				},
			},
			entitiesToRemove: []entityToRemove{
				{
					entityID: 0,
					want: map[reflect.Type]any{
						reflect.TypeFor[int](): helpers.New(1),
					},
					wantErr: false,
				},
				{
					entityID: 0,
					want: map[reflect.Type]any{
						reflect.TypeFor[int](): helpers.New(1),
					},
					wantErr: true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, err := NewArchetype(tt.componentTypes, tt.signatureMask)
			require.NoError(t, err)

			// Add entities
			for _, entityToAdd := range tt.entitiesToAdd {
				gotErr := a.AddEntity(entityToAdd.entityID, entityToAdd.componentsData)
				require.NoError(t, gotErr)
			}

			for _, entityToRemove := range tt.entitiesToRemove {
				got, gotErr := a.RemoveEntity(entityToRemove.entityID)
				if entityToRemove.wantErr {
					require.Error(t, gotErr)
					return
				}
				require.NoError(t, gotErr)
				require.Equal(t, entityToRemove.want, got)
				require.False(t, slices.ContainsFunc(a.entities, func(e EntityID) bool {
					return e == entityToRemove.entityID
				}))
				for componentType := range entityToRemove.want {
					_, exists := a.GetComponentPtr(entityToRemove.entityID, componentType)
					require.False(t, exists)
				}
			}
		})
	}
}
