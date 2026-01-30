package netcode

type PeerID = uint64

type NetworkRole int

const (
	RoleNone NetworkRole = iota
	RoleServer
	RoleClient
	RoleHost
)

type EventType int

const (
	EventData EventType = iota
	EventConnect
	EventDisconnect
)

// Transport defines how bytes are sent/received (UDP/TCP/WebSockets).
type Transport interface {
	Send(peerID PeerID, data []byte) error
	Poll() []NetworkEvent
	ClientID() PeerID // Local ID
}

type NetworkEvent struct {
	Type   EventType
	Sender PeerID
	Data   []byte
}

// NetworkManager handles the global network state.
type NetworkManager struct {
	role      NetworkRole
	transport Transport
	nextNetID PeerID // For generating NetworkIDs (Server only)
}

func NewNetworkManager(transport Transport, role NetworkRole) *NetworkManager {
	return &NetworkManager{
		role:      role,
		transport: transport,
		nextNetID: 1,
	}
}

func (nm *NetworkManager) IsServer() bool {
	return nm.role == RoleServer || nm.role == RoleHost
}

func (nm *NetworkManager) IsClient() bool {
	return nm.role == RoleClient || nm.role == RoleHost
}

func (nm *NetworkManager) IsHost() bool {
	return nm.role == RoleHost
}

func (nm *NetworkManager) PeerID() PeerID {
	return nm.transport.ClientID()
}
