package message

import (
	"fmt"
)

// Type is the ID for the message type describing which payload is in the payload field of the message.
type Type int32

const (
	TypeConnect                  Type = 0   // TypeConnect describes a message that is sent when the user connects on a TCP level.
	TypeDisconnect               Type = 1   // TypeDisconnect describes a message that is sent when the user disconnects on a TCP level.
	TypeAuthPassword             Type = 100 // TypeAuthPassword describes a message that is sent when the user submits a username and password.
	TypeAuthPasswordSuccessful   Type = 101 // TypeAuthPasswordSuccessful describes a message that is sent when the submitted username and password were valid.
	TypeAuthPasswordFailed       Type = 102 // TypeAuthPasswordFailed describes a message that is sent when the submitted username and password were invalid.
	TypeAuthPasswordBackendError Type = 103 // TypeAuthPasswordBackendError describes a message that is sent when the auth server failed to respond to a request with username and password
	TypeAuthPubKey               Type = 104 // TypeAuthPubKey describes a message that is sent when the user submits a username and public key.
	TypeAuthPubKeySuccessful     Type = 105 // TypeAuthPubKeySuccessful describes a message that is sent when the submitted username and public key were invalid.
	TypeAuthPubKeyFailed         Type = 106 // TypeAuthPubKeyFailed describes a message that is sent when the submitted username and public key were invalid.
	TypeAuthPubKeyBackendError   Type = 107 // TypeAuthPubKeyBackendError describes a message that is sent when the auth server failed to respond with username and password.

	TypeAuthKeyboardInteractiveChallenge    Type = 108 // TypeAuthKeyboardInteractiveChallenge is a message that indicates that a keyboard-interactive challenge has been sent to the user. Multiple challenge-response interactions can take place.
	TypeAuthKeyboardInteractiveAnswer       Type = 109 // TypeAuthKeyboardInteractiveAnswer is a message that indicates a response to a keyboard-interactive challenge.
	TypeAuthKeyboardInteractiveFailed       Type = 110 // TypeAuthKeyboardInteractiveFailed indicates that a keyboard-interactive authentication process has failed.
	TypeAuthKeyboardInteractiveBackendError Type = 111 // TypeAuthKeyboardInteractiveBackendError indicates an error in the authentication backend during a keyboard-interactive authentication.

	TypeHandshakeFailed      Type = 198 // TypeHandshakeFailed indicates that the handshake has failed.
	TypeHandshakeSuccessful  Type = 199 // TypeHandshakeSuccessful indicates that the handshake and authentication was successful.
	TypeGlobalRequestUnknown Type = 200 // TypeGlobalRequestUnknown describes a message when a global (non-channel) request was sent that was not recognized.

	TypeNewChannel           Type = 300 // TypeNewChannel describes a message that indicates a new channel request.
	TypeNewChannelSuccessful Type = 301 // TypeNewChannelSuccessful describes a message when the new channel request was successful.
	TypeNewChannelFailed     Type = 302 // TypeNewChannelFailed describes a message when the channel request failed for the reason indicated.

	TypeChannelRequestUnknownType  Type = 400 // TypeChannelRequestUnknownType describes an in-channel request from the user that is not supported.
	TypeChannelRequestDecodeFailed Type = 401 // TypeChannelRequestDecodeFailed describes an in-channel request from the user that is supported but the payload could not be decoded.
	TypeChannelRequestSetEnv       Type = 402 // TypeChannelRequestSetEnv describes an in-channel request to set an environment variable.
	TypeChannelRequestExec         Type = 403 // TypeChannelRequestExec describes an in-channel request to run a program.
	TypeChannelRequestPty          Type = 404 // TypeChannelRequestPty describes an in-channel request to create an interactive terminal

	TypeChannelRequestShell  Type = 405 // TypeChannelRequestShell describes an in-channel request to start a shell.
	TypeChannelRequestSignal Type = 406 // TypeChannelRequestSignal describes an in-channel request to send a signal to the currently running program.

	TypeChannelRequestSubsystem Type = 407 // TypeChannelRequestSubsystem describes an in-channel request to start a well-known subsystem (e.g. SFTP).
	TypeChannelRequestWindow    Type = 408 // TypeChannelRequestWindow describes an in-channel request to resize the current interactive terminal.

	TypeWriteClose Type = 496 // TypeWriteClose indicates that the channel was closed for writing from the server side.
	TypeClose      Type = 497 // TypeClose indicates that the channel was closed.
	TypeExitSignal Type = 498 // TypeExitSignal describes the signal that caused a program to terminate abnormally.
	TypeExit       Type = 499 // TypeExit describes a message that is sent when the program exited. The payload contains the exit status.

	TypeIO            Type = 500 // TypeIO describes the testdata transferred to and from the currently running program on the terminal.
	TypeRequestFailed Type = 501 // TypeRequestFailed describes that a request has failed.
)

var typeToID = map[Type]string{
	TypeConnect:    "connect",
	TypeDisconnect: "disconnect",

	TypeAuthPassword:             "auth_password",
	TypeAuthPasswordSuccessful:   "auth_password_successful",
	TypeAuthPasswordFailed:       "auth_password_failed",
	TypeAuthPasswordBackendError: "auth_password_backend_error",

	TypeAuthPubKey:             "auth_pubkey",
	TypeAuthPubKeySuccessful:   "auth_pubkey_successful",
	TypeAuthPubKeyFailed:       "auth_pubkey_failed",
	TypeAuthPubKeyBackendError: "auth_pubkey_backend_error",

	TypeAuthKeyboardInteractiveChallenge:    "auth_keyboard_interactive_challenge",
	TypeAuthKeyboardInteractiveAnswer:       "auth_keyboard_interactive_answer",
	TypeAuthKeyboardInteractiveFailed:       "auth_keyboard_interactive_failed",
	TypeAuthKeyboardInteractiveBackendError: "auth_keyboard_interactive_backend_error",

	TypeHandshakeFailed:     "handshake_failed",
	TypeHandshakeSuccessful: "handshake_successful",

	TypeGlobalRequestUnknown: "global_request_unknown",
	TypeNewChannel:           "new_channel",
	TypeNewChannelSuccessful: "new_channel_successful",
	TypeNewChannelFailed:     "new_channel_failed",

	TypeChannelRequestUnknownType:  "channel_request_unknown",
	TypeChannelRequestDecodeFailed: "channel_request_decode_failed",
	TypeChannelRequestSetEnv:       "setenv",
	TypeChannelRequestExec:         "exec",
	TypeChannelRequestPty:          "pty",
	TypeChannelRequestShell:        "shell",
	TypeChannelRequestSignal:       "signal",
	TypeChannelRequestSubsystem:    "subsystem",
	TypeChannelRequestWindow:       "window",
	TypeWriteClose:                 "close_write",
	TypeClose:                      "close",
	TypeExit:                       "exit",
	TypeExitSignal:                 "exit_signal",

	TypeIO:            "io",
	TypeRequestFailed: "request_failed",
}

var typeToName = map[Type]string{
	TypeConnect:    "Connect",
	TypeDisconnect: "Disconnect",

	TypeAuthPassword:             "Password authentication",
	TypeAuthPasswordSuccessful:   "Password authentication successful",
	TypeAuthPasswordFailed:       "Password authentication failed",
	TypeAuthPasswordBackendError: "Password authentication backend error",

	TypeAuthPubKey:             "Public key authentication",
	TypeAuthPubKeySuccessful:   "Public key authentication successful",
	TypeAuthPubKeyFailed:       "Public key authentication failed",
	TypeAuthPubKeyBackendError: "Public key authentication backend error",

	TypeAuthKeyboardInteractiveChallenge:    "Keyboard-interactive authentication challenge",
	TypeAuthKeyboardInteractiveAnswer:       "Keyboard-interactive authentication answer",
	TypeAuthKeyboardInteractiveFailed:       "Keyboard-interactive authentication failed",
	TypeAuthKeyboardInteractiveBackendError: "Keyboard-interactive authentication backend error",

	TypeHandshakeFailed:     "Handshake failed",
	TypeHandshakeSuccessful: "Handshake successful",

	TypeGlobalRequestUnknown: "Unknown global request",
	TypeNewChannel:           "New channel request",
	TypeNewChannelSuccessful: "New channel successful",
	TypeNewChannelFailed:     "New channel failed",

	TypeChannelRequestUnknownType:  "Unknown channel request",
	TypeChannelRequestDecodeFailed: "Failed to decode channel request",
	TypeChannelRequestSetEnv:       "Set environment variable",
	TypeChannelRequestExec:         "Execute program",
	TypeChannelRequestPty:          "Request interactive terminal",
	TypeChannelRequestShell:        "Run shell",
	TypeChannelRequestSignal:       "Send signal to running process",
	TypeChannelRequestSubsystem:    "Request subsystem",
	TypeChannelRequestWindow:       "Change window size",
	TypeWriteClose:                 "Close channel for writing",
	TypeClose:                      "Close channel",
	TypeExit:                       "Program exited",
	TypeExitSignal:                 "Program exited with signal",

	TypeIO:            "I/O",
	TypeRequestFailed: "Request failed",
}

var messageTypeToPayload = map[Type]Payload{
	TypeConnect:    PayloadConnect{},
	TypeDisconnect: nil,

	TypeAuthPassword:                        PayloadAuthPassword{},
	TypeAuthPasswordSuccessful:              PayloadAuthPassword{},
	TypeAuthPasswordFailed:                  PayloadAuthPassword{},
	TypeAuthPasswordBackendError:            PayloadAuthPasswordBackendError{},
	TypeAuthKeyboardInteractiveChallenge:    PayloadAuthKeyboardInteractiveChallenge{},
	TypeAuthKeyboardInteractiveAnswer:       PayloadAuthKeyboardInteractiveAnswer{},
	TypeAuthKeyboardInteractiveFailed:       PayloadAuthKeyboardInteractiveFailed{},
	TypeAuthKeyboardInteractiveBackendError: PayloadAuthKeyboardInteractiveBackendError{},
	TypeHandshakeFailed:                     PayloadHandshakeFailed{},
	TypeHandshakeSuccessful:                 PayloadHandshakeSuccessful{},

	TypeAuthPubKey:             PayloadAuthPubKey{},
	TypeAuthPubKeySuccessful:   PayloadAuthPubKey{},
	TypeAuthPubKeyFailed:       PayloadAuthPubKey{},
	TypeAuthPubKeyBackendError: PayloadAuthPubKeyBackendError{},

	TypeGlobalRequestUnknown: PayloadGlobalRequestUnknown{},
	TypeNewChannel:           PayloadNewChannel{},
	TypeNewChannelSuccessful: PayloadNewChannelSuccessful{},
	TypeNewChannelFailed:     PayloadNewChannelFailed{},

	TypeChannelRequestUnknownType:  PayloadChannelRequestUnknownType{},
	TypeChannelRequestDecodeFailed: PayloadChannelRequestDecodeFailed{},
	TypeChannelRequestSetEnv:       PayloadChannelRequestSetEnv{},
	TypeChannelRequestExec:         PayloadChannelRequestExec{},
	TypeChannelRequestPty:          PayloadChannelRequestPty{},
	TypeChannelRequestShell:        PayloadChannelRequestShell{},
	TypeChannelRequestSignal:       PayloadChannelRequestSignal{},
	TypeChannelRequestSubsystem:    PayloadChannelRequestSubsystem{},
	TypeChannelRequestWindow:       PayloadChannelRequestWindow{},
	TypeIO:                         PayloadIO{},
	TypeRequestFailed:              PayloadRequestFailed{},
	TypeExit:                       PayloadExit{},
	TypeExitSignal:                 PayloadExitSignal{},

	TypeClose:      nil,
	TypeWriteClose: nil,
}

// ListTypes returns all defined types.
func ListTypes() []Type {
	var keys []Type
	for key := range typeToID {
		keys = append(keys, key)
	}
	return keys
}

// ID converts the numeric message type to a string representation for human consumption.
func (messageType Type) ID() string {
	if val, ok := typeToID[messageType]; ok {
		return val
	}
	return "invalid"
}

// Name converts the numeric message type to a string representation for human consumption.
func (messageType Type) Name() string {
	if val, ok := typeToName[messageType]; ok {
		return val
	}
	return "invalid"
}

// Code returns a numeric code for this message type.
func (messageType Type) Code() int32 {
	return int32(messageType)
}

// Payload returns a typed struct for a payload. May be nil if the payload is empty.
func (messageType Type) Payload() (Payload, error) {
	payload, ok := messageTypeToPayload[messageType]
	if !ok {
		return nil, fmt.Errorf("invalid message type: %d", messageType)
	}
	return payload, nil
}
