package agentprotocol

type Connection interface {
	Read(p []byte) (n int, err error)
	Write(data []byte) (n int, err error)
	Close() error
	CloseImmediately() error
	Accept() error
	Reject() error
	Details() NewConnectionPayload
}
