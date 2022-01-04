package sshserver

import (
	"context"
	"fmt"

	auth2 "github.com/containerssh/libcontainerssh/auth"
	"github.com/containerssh/libcontainerssh/internal/auth"
)

// AbstractNetworkConnectionHandler is an empty implementation for the NetworkConnectionHandler interface.
type AbstractNetworkConnectionHandler struct {
}

// OnAuthPassword is called when a user attempts a password authentication. The implementation must always supply
//                AuthResponse and may supply error as a reason description.
func (a *AbstractNetworkConnectionHandler) OnAuthPassword(_ string, _ []byte, _ string) (response AuthResponse, metadata *auth2.ConnectionMetadata, reason error) {
	return AuthResponseUnavailable, nil, nil
}

// OnAuthPassword is called when a user attempts a pubkey authentication. The implementation must always supply
//                AuthResponse and may supply error as a reason description. The pubKey parameter is an SSH key in
//               the form of "ssh-rsa KEY HERE".
func (a *AbstractNetworkConnectionHandler) OnAuthPubKey(_ string, _ string, _ string) (response AuthResponse, metadata *auth2.ConnectionMetadata, reason error) {
	return AuthResponseUnavailable, nil, nil
}

// OnAuthKeyboardInteractive is a callback for interactive authentication. The implementer will be passed a callback
// function that can be used to issue challenges to the user. These challenges can, but do not have to contain
// questions.
func (a *AbstractNetworkConnectionHandler) OnAuthKeyboardInteractive(
	_ string,
	_ func(
		instruction string,
		questions KeyboardInteractiveQuestions,
	) (answers KeyboardInteractiveAnswers, err error),
	_ string,
) (response AuthResponse, metadata *auth2.ConnectionMetadata, reason error) {
	return AuthResponseUnavailable, nil, nil
}

func (a *AbstractNetworkConnectionHandler) OnAuthGSSAPI() auth.GSSAPIServer {
	return nil
}

// OnHandshakeFailed is called when the SSH handshake failed. This method is also called after an authentication
//                   failure. After this method is the connection will be closed and the OnDisconnect method will be
//                   called.
func (a *AbstractNetworkConnectionHandler) OnHandshakeFailed(_ error) {

}

// OnHandshakeSuccess is called when the SSH handshake was successful. It returns connection to process
//                    requests, or failureReason to indicate that a backend error has happened. In this case, the
//                    connection will be closed and OnDisconnect will be called.
func (a *AbstractNetworkConnectionHandler) OnHandshakeSuccess(_ string) (
	connection SSHConnectionHandler, failureReason error,
) {
	return nil, fmt.Errorf("not implemented")
}

// OnDisconnect is called when the network connection is closed.
func (a *AbstractNetworkConnectionHandler) OnDisconnect() {
}

// OnShutdown is called when a shutdown of the SSH server is desired. The shutdownContext is passed as a deadline
//            for the shutdown, after which the server should abort all running connections and return as fast as
//            possible.
func (a *AbstractNetworkConnectionHandler) OnShutdown(_ context.Context) {}
