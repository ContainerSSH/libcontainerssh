package auth

import (
	"go.containerssh.io/libcontainerssh/metadata"
)

// PasswordAuthRequest is an authentication request for password authentication.
//
// swagger:model PasswordAuthRequest
type PasswordAuthRequest struct {
	// swagger:allOf
	metadata.ConnectionAuthPendingMetadata `json:",inline"`

	// Password the user provided for authentication.
	//
	// required: true
	// in: body
	// swagger:strfmt Base64
	Password string `json:"passwordBase64"`
}

// PublicKeyAuthRequest is an authentication request for public key authentication.
//
// swagger:model PublicKeyAuthRequest
type PublicKeyAuthRequest struct {
	// swagger:allOf
	metadata.ConnectionAuthPendingMetadata `json:",inline"`

	// in: body
	// required: true
	PublicKey `json:",inline"`
}

// AuthorizationRequest is the authorization request used after some
// authentication methods (e.g. kerberos) to determine whether users are
// allowed to access the service
//
// swagger:model AuthorizationRequest
type AuthorizationRequest struct {
	// swagger:allOf
	metadata.ConnectionAuthenticatedMetadata `json:",inline"`
}

// ResponseBody is a response to authentication requests.
//
// swagger:model AuthResponseBody
type ResponseBody struct {
	metadata.DynamicMetadata `json:",inline"`

	// AuthenticatedUsername contains the username that was actually verified. This may differ from LoginUsername when,
	// for example OAuth2 or Kerberos authentication is used. This field is empty until the authentication phase is
	// completed.
	//
	// required: false
	// in: body
	// example: systemusername
	AuthenticatedUsername string `json:"authenticatedUsername,omitempty"`

	// Success indicates if the authentication was successful.
	//
	// required: true
	// in: body
	Success bool `json:"success"`
}

// Response is the full HTTP authentication response.
//
// swagger:response AuthResponse
type Response struct {
	// The response body
	//
	// in: body
	// required: true
	Body ResponseBody
}
