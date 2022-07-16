package agentprotocol

import "io"

type ForwardCtx interface {
	NewConnectionTCP(
		connectedAddress string,
		connectedPort uint32,
		origAddress string,
		origPort uint32,
		closeFunc func() error,
	) (io.ReadWriteCloser, error)
	NewConnectionUnix(
		path string,
		closeFunc func() error,
	) (io.ReadWriteCloser, error)

	StartClient() (connectionType uint32, setupPacket SetupPacket, connChan chan Connection, err error)
	StartServerForward() (chan Connection, error)
	StartX11ForwardClient(
		singleConnection bool,
		screen string,
		authProtocol string,
		authCookie string,
	) (chan Connection, error)
	StartReverseForwardClient(bindHost string, bindPort uint32, singleConnection bool) (chan Connection, error)
	StartReverseForwardClientUnix(path string, singleConnection bool) (chan Connection, error)
	NoMoreConnections() error
	WaitFinish()
	Kill()
}
