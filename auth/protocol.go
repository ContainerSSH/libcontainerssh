package auth

import "time"

// PasswordAuthRequest is an authentication request for password authentication.
//
// swagger:model PasswordAuthRequest
type PasswordAuthRequest struct {
	// Username is the username provided for authentication.
	//
	// required: true
	Username string `json:"username"`

	// RemoteAddress is the IP address of the user trying to authenticate.
	//
	// required: true
	RemoteAddress string `json:"remoteAddress"`

	// ConnectionID is an opaque ID to identify the SSH connection in question.
	//
	// required: true
	ConnectionID string `json:"connectionId"`

	// SessionID is a deprecated alias for ConnectionID and will be removed in the future.
	//
	// required: true
	SessionID string `json:"sessionId"`

	// Password the user provided for authentication.
	//
	// required: true
	Password string `json:"passwordBase64"`
}

// PublicKeyAuthRequest is an authentication request for public key authentication.
//
// swagger:model PublicKeyAuthRequest
type PublicKeyAuthRequest struct {
	// Username is the username provided for authentication.
	//
	// required: true
	Username string `json:"username"`

	// RemoteAddress is the IP address of the user trying to authenticate.
	//
	// required: true
	RemoteAddress string `json:"remoteAddress"`

	// ConnectionID is an opaque ID to identify the SSH connection in question.
	//
	// required: true
	ConnectionID string `json:"connectionId"`

	// SessionID is a deprecated alias for ConnectionID and will be removed in the future.
	//
	// required: true
	SessionID string `json:"sessionId"`

	// PublicKey is the key in the authorized key format.
	//
	// required: true
	PublicKey string `json:"publicKey"`

	// CACertificate contains information about the SSH certificate presented by a connecting client. This certificate
	// is not an SSL/TLS/x509 certificate and has a much simpler structure. However, this can be used to verify if the
	// connecting client belongs to an organization.
	//
	// required: false
	CACertificate *CACertificate `json:"caCertificate,omitempty"`
}

// ResponseBody is a response to authentication requests.
//
// swagger:model AuthResponseBody
type ResponseBody struct {
	// Success indicates if the authentication was successful.
	//
	// required: true
	Success bool `json:"success"`

	// Metadata is a set of key-value pairs that can be returned and either consumed by the configuration server or
	// exposed in the backend as environment variables.
	//
	// required: false
	Metadata map[string]string `json:"metadata,omitempty"`
}

// Response is the full HTTP authentication response.
//
// swagger:response AuthResponse
type Response struct {
	// The response body
	//
	// in: body
	ResponseBody
}

// CACertificate contains information about the SSH certificate presented by a connecting client. This certificate
// is not an SSL/TLS/x509 certificate and has a much simpler structure. However, this can be used to verify if the
// connecting client belongs to an organization.
//
// swagger:model CACertificate
type CACertificate struct {
	// PublicKey contains the public key of the CA signing the public key presented in the OpenSSH authorized key
	// format.
	PublicKey string `json:"key"`
	// KeyID contains an identifier for the key.
	KeyID string `json:"keyID"`
	// ValidPrincipals contains a list of principals for which this CA certificate is valid.
	ValidPrincipals []string `json:"validPrincipals"`
	// ValidAfter contains the time after which this certificate is valid.
	ValidAfter time.Time `json:"validAfter"`
	// ValidBefore contains the time when this certificate expires.
	ValidBefore time.Time `json:"validBefore"`
}
