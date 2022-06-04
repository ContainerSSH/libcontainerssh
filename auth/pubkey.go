package auth

// PublicKey contains the details of a public key provided during authentication.
type PublicKey struct {
	// PublicKey is the key in the authorized key format.
	//
	// required: true
	// example: ssh-rsa ...
	PublicKey string `json:"publicKey"`
}
