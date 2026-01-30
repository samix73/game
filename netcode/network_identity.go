package netcode

import ecs "github.com/samix73/ebiten-ecs"

func init() {
	ecs.RegisterComponent[NetworkIdentity]()
}

// NetworkIdentity identifies an entity across the network.
type NetworkIdentity struct {
	NetworkID PeerID
	OwnerID   PeerID // The ClientID that owns this entity (0 = Server)
	Spawned   bool
	IsPlayer  bool
}

func (n *NetworkIdentity) Init() {
	n.NetworkID = 0
	n.OwnerID = 0
	n.Spawned = false
	n.IsPlayer = false
}

func (n *NetworkIdentity) Reset() {
	n.NetworkID = 0
	n.OwnerID = 0
	n.Spawned = false
	n.IsPlayer = false
}

func (n *NetworkIdentity) IsOwner(localClientID PeerID) bool {
	return n.OwnerID == localClientID
}
