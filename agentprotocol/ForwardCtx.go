package agentprotocol

import "io"

type ForwardCtx interface {
	// NewConnectionTCP requests the other side to connect to a specified address/host and port combination and forward all data from the returned ReadWriteCloser to it.
	//
	// connectedAddress is the address that the connection requested to connect to, is it used by the receiving side to initiate the connection to the desired address.
	// connectedPort is the port that the connection requested to connect to, is it used by the receiving side to initiate the connection to the desired port.
	// origAddress is the originator address of the connection. It can be used by the receiving side to decide whether to accept this connection.
	// origAddress is the originator port of the connection. It can be used by the receiving side to decide whether to accept this connection.
	// closeFunc is a callback function that is called when the connection is called to perform cleanup of the backing connection.
	NewConnectionTCP(
		connectedAddress string,
		connectedPort uint32,
		origAddress string,
		origPort uint32,
		closeFunc func() error,
	) (io.ReadWriteCloser, error)
	// NewConnectionUnix requests the other side to connect to a specified unix and forward all data from the returned ReadWriteCloser to it.
	//
	// path is the path of the unix socket to connect to.
	// closeFunc is a callback function that is called when the connection is called to perform cleanup of the backing connection.
	NewConnectionUnix(
		path string,
		closeFunc func() error,
	) (io.ReadWriteCloser, error)

	// StartServer initializes the ForwardCtx in server mode which waits for information from the other side about the function it needs to perform. It returns the connection type the other side requests and additional information in the setupPacket. Additionally, a Connection channel, connChan, is returned that provides connection requests from the other side of the connection.
	StartServer() (connectionType uint32, setupPacket SetupPacket, connChan chan Connection, err error)

	// StartClientForward initializes the ForwardCtx in client mode and informs the server that the client is going to be the connection requestor (Direct Forward). A connection channel is returned that informs the server of connection requests by the client however in this mode it is a assumed that the server sends no connection request so the sane behaviour is to reject all connections. In this mode, the client can start new connections on the server using the NewConnection* function family.
	StartClientForward() (chan Connection, error)

	// StartX11ForwardClient initializes the ForwardCtx in client mode and informs the server to start an X11 server and forward all X11 connections to the client.
	//
	// singleConnection is the X11 singleConnection parameter that requests to only accept the first connection (X11 window) and no more.
	// screen is the X11 screen number.
	// authProtocol is the X11 auth protocol.
	// authCookie is the X11 auth cookie.
	StartX11ForwardClient(
		singleConnection bool,
		screen string,
		authProtocol string,
		authCookie string,
	) (chan Connection, error)

	// StartReverseForwardClient initializes the ForwardCtx in client mode and informs the server to start listening for connections on the requested host and port. Once a connection is received a new connection is created and sent through the Connection channel.
	//
	// bindHost is the host to listen on connections on.
	// bindPort is the port to listen on connections on.
	// singleConnection is a flag that requests to stop listening for new connections after the first one.
	StartReverseForwardClient(bindHost string, bindPort uint32, singleConnection bool) (chan Connection, error)

	// StartReverseForwardClient initializes the ForwardCtx in client mode and informs the server to start listening for connections on the requested unix socket. Once a connection is received a new connection is created and sent through the Connection channel.
	//
	// path is the path to the unix socket to listen on.
	// singleConnection is a flag that requests to stop listening for new connections after the first one.
	StartReverseForwardClientUnix(path string, singleConnection bool) (chan Connection, error)

	// NoMoreConnections informs the other side that it should not accept any more connection requests from it. It is used as a forwarding security feature in cases where it's clear there will only be one connection.
	NoMoreConnections() error

	// WaitFinish blocks until NoMoreConnections has been received and all active connections have been closed.
	WaitFinish()

	// Kill closes the ForwardCtx immediately and terminates all connections.
	Kill()
}
