package agentprotocol

type Connection interface {
	// Read reads data from the connection. The blocking nature of this call depends on the underlying communication medium
	Read(p []byte) (n int, err error)
	// Read reads data from the connection. The blocking nature of this call depends on the underlying communication medium
	Write(data []byte) (n int, err error)
	// Close requests to close an active connection
	Close() error
	// CloseImmediately closes an active connection without waiting for the other side to acknowledge
	CloseImmediately() error
	// Accept accepts a pending connection request. It must be called before any Read/Write functions can be called on the connection
	Accept() error
	// Reject rejects a pending connection request and closes the connection.
	Reject() error
	// Details returns the details of a connection request. It can be called to gain more information about a connection before an Accept/Reject action is made
	Details() NewConnectionPayload
}
