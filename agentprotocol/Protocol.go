package agentprotocol

const (
	CONNECTION_TYPE_X11 = iota
	CONNECTION_TYPE_PORT_FORWARD
	CONNECTION_TYPE_PORT_DIAL
	CONNECTION_TYPE_SOCKET_FORWARD
	CONNECTION_TYPE_SOCKET_DIAL
)

const (
	PROTOCOL_TCP  string = "tcp"
	PROTOCOL_UNIX string = "unix"
)

const (
	PACKET_SETUP = iota
	PACKET_SUCCESS
	PACKET_ERROR
	PACKET_DATA
	PACKET_NEW_CONNECTION
	PACKET_CLOSE_CONNECTION
	PACKET_NO_MORE_CONNECTIONS
)

type SetupPacket struct {
	ConnectionType uint32
	BindHost       string
	BindPort       uint32
	Protocol       string

	Screen           string
	SingleConnection bool
	AuthProtocol     string
	AuthCookie       string
}

type NewConnectionPayload struct {
	Protocol string

	ConnectedAddress  string
	ConnectedPort     uint32
	OriginatorAddress string
	OriginatorPort    uint32
}

type Packet struct {
	Type         int
	ConnectionID uint64
	Payload      []byte
}
