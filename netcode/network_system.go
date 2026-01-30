package netcode

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"reflect"

	ecs "github.com/samix73/ebiten-ecs"
	"github.com/samix73/game/helpers"
)

// Packet structure
type PacketHeader struct {
	MsgType     byte
	NetworkID   PeerID
	ComponentID ecs.ComponentID
}

const (
	MsgVarUpdate byte = 1
	MsgSpawn     byte = 2
)

// NetworkSystem handles replication.
type NetworkSystem struct {
	*ecs.BaseSystem
	Manager *NetworkManager
}

func NewNetworkSystem(priority int, manager *NetworkManager) *NetworkSystem {
	return &NetworkSystem{
		BaseSystem: ecs.NewBaseSystem(priority),
		Manager:    manager,
	}
}

func (s *NetworkSystem) Update() error {
	// 1. Process Incoming Packets
	if err := s.processIncomingMessages(); err != nil {
		return fmt.Errorf("failed to process incoming messages: %w", err)
	}

	// 2. Replicate Outgoing State (Server -> Clients)
	if s.Manager.IsServer() {
		if err := s.replicateVariables(); err != nil {
			return fmt.Errorf("failed to replicate variables: %w", err)
		}
	}

	return nil
}

func (s *NetworkSystem) processIncomingMessages() error {
	for _, event := range s.Manager.transport.Poll() {
		if event.Type == EventData {
			if err := s.handlePacket(event.Data); err != nil {
				return fmt.Errorf("failed to handle packet: %w", err)
			}
		}
	}

	return nil
}

func (s *NetworkSystem) handlePacket(data []byte) error {
	var header PacketHeader
	if err := gob.NewDecoder(bytes.NewBuffer(data)).Decode(&header); err != nil {
		return fmt.Errorf("failed to decode packet header: %w", err)
	}

	if header.MsgType != MsgVarUpdate {
		return fmt.Errorf("unknown message type: %d", header.MsgType)
	}

	// Find entity by NetworkID (This requires a lookup map in a real scenario, linear scan used for brevity)
	targetEntity, ok := s.findEntityByNetID(header.NetworkID)
	if !ok {
		return fmt.Errorf("entity not found for network id: %d", header.NetworkID)
	}

	// Apply update to the specific component
	if err := s.applyComponentUpdate(targetEntity, header.ComponentID, data); err != nil {
		return fmt.Errorf("failed to apply component update: %w", err)
	}

	return nil
}

func (s *NetworkSystem) replicateVariables() error {
	em := s.EntityManager()

	// Query all entities with a NetworkIdentity
	entities := ecs.QueryWith(em, func(nid *NetworkIdentity) bool {
		return nid.NetworkID != 0
	})
	for range entities {
		// Iterate ALL components on this entity to find NetworkVariables
		// Note: In a production generic ECS, you might optimize this by registering "NetworkComponents" separately.
		// Here, we inspect the components available in the archetype.

		// This uses internal knowledge or you must track which components to scan.
		// For this example, let's assume we scan a specific component type provided by the user,
		// or we iterate a known list of component types.
		// Since ebiten-ecs doesn't expose "GetAllComponents", we rely on the user registering components to sync.
		// LIMITATION: You must manually check components you want to sync here, or create a 'Replicable' interface.
	}

	return nil
}

// SyncComponent is a helper to manually call in your systems or update loop
// It scans the struct for NetworkVariable fields.
func (s *NetworkSystem) SyncComponent(entityID ecs.EntityID, component any, compID ecs.ComponentID) {
	if !s.Manager.IsServer() {
		return
	}

	val := reflect.ValueOf(component)
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}

	netIdentity, _ := ecs.GetComponent[NetworkIdentity](s.EntityManager(), entityID)
	if netIdentity == nil || netIdentity.NetworkID == 0 {
		return
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		// Check if field implements NetworkVariableInterface
		if field.CanAddr() {
			if netVar, ok := field.Addr().Interface().(NetworkVariableInterface); ok {
				if netVar.IsDirty() {
					// Serialize and Send
					payload, _ := netVar.Serialize()

					var buf bytes.Buffer
					enc := gob.NewEncoder(&buf)

					// Header
					header := PacketHeader{
						MsgType:     MsgVarUpdate,
						NetworkID:   netIdentity.NetworkID,
						ComponentID: compID,
					}
					enc.Encode(header)
					buf.Write(payload)

					// Broadcast to all clients (Implementation specific to Transport)
					// s.Manager.transport.Broadcast(buf.Bytes())

					netVar.ResetDirty()
				}
			}
		}
	}
}

func (s *NetworkSystem) applyComponentUpdate(entityID ecs.EntityID, compID ecs.ComponentID, data []byte) error {
	// This part is tricky in Go generics without a central type registry that returns 'any'.
	// In ebiten-ecs, we have 'componentsPools' in registry.go but it's internal.
	// You would need to expose a way to Get component by ID dynamically.

	return nil
}

func (s *NetworkSystem) findEntityByNetID(netID uint64) (ecs.EntityID, bool) {
	// In production, maintain a map[uint64]EntityID cache in this System.
	// Linear scan for example:
	em := s.EntityManager()
	entities := ecs.QueryWith(em, func(nid *NetworkIdentity) bool {
		return nid.NetworkID == netID
	})

	return helpers.First(entities)
}

func (s *NetworkSystem) Teardown() {}
