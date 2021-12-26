package message

import "time"

// CACertificate is an SSH certificate presented by a client to verify their key against a CA.
type CACertificate struct {
	// PublicKey contains the public key of the CA signing the public key presented in the OpenSSH authorized key
	// format.
	PublicKey string `json:"key"`
	// KeyID contains an identifier for the key.
	KeyID string `json:"keyID"`
	// ValidPrincipals contains a list of principals for which this CA certificate is valid.
	ValidPrincipals []string `json:"validPrincipals"`
	// ValidAfter contains the time after which this certificate is valid. This may be empty.
	ValidAfter time.Time `json:"validAfter"`
	// ValidBefore contains the time when this certificate expires. This may be empty.
	ValidBefore time.Time `json:"validBefore"`
}

// Equals compares the CACertificate record.
func (c *CACertificate) Equals(cert *CACertificate) bool {
	if c == nil {
		return cert == nil
	}
	if cert == nil {
		return false
	}
	if c.PublicKey != cert.PublicKey {
		return false
	}
	if c.KeyID != cert.KeyID {
		return false
	}
	if len(c.ValidPrincipals) != len(cert.ValidPrincipals) {
		return false
	}
	for _, validPrincipal := range c.ValidPrincipals {
		found := false
		for _, otherPrincipal := range cert.ValidPrincipals {
			if otherPrincipal == validPrincipal {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return c.ValidAfter == cert.ValidAfter && c.ValidBefore == cert.ValidBefore
}
