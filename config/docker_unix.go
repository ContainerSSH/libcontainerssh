//go:build linux || freebsd || openbsd || darwin
// +build linux freebsd openbsd darwin

package config

// DockerConnectionConfig configures how to connect to dockerd.
type DockerConnectionConfig struct {
	// Host is the docker connect URL.
	Host string `json:"host,omitempty" yaml:"host,omitempty" default:"unix:///var/run/docker.sock"`
	// CaCert is the CA certificate for Docker connection embedded in the configuration in PEM format.
	CaCert string `json:"cacert,omitempty" yaml:"cacert,omitempty"`
	// Cert is the client certificate in PEM format embedded in the configuration.
	Cert string `json:"cert,omitempty" yaml:"cert,omitempty"`
	// Key is the client key in PEM format embedded in the configuration.
	Key string `json:",omitempty" yaml:"key,omitempty"`
}
