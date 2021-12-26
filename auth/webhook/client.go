package webhook

import (
	"net"

	protocol "github.com/containerssh/libcontainerssh/auth"
	"github.com/containerssh/libcontainerssh/config"
	"github.com/containerssh/libcontainerssh/internal/auth"
	"github.com/containerssh/libcontainerssh/internal/geoip/dummy"
	"github.com/containerssh/libcontainerssh/internal/metrics"
	"github.com/containerssh/libcontainerssh/log"
)

type Client interface {
	// Password authenticates with a password from the client. It returns a bool if the authentication as successful
	// or not. If an error happened while contacting the authentication server it will return an error.
	Password(
		username string,
		password []byte,
		connectionID string,
		remoteAddr net.IP,
	) AuthenticationContext

	// PubKey authenticates with a public key from the client. It returns a bool if the authentication as successful
	// or not. If an error happened while contacting the authentication server it will return an error.
	//
	// The parameters are as follows:
	//
	// - username is the username provided by the connecting client.
	// - pubKey is the public key offered by the connecting client. The client may offer multiple keys which will be
	// presented by calling this function multiple times.
	// - connectionID is an opaque random string representing this SSH connection across multiple webhooks and logs.
	// - remoteAddr is the IP address of the connecting client.
	// - caPubKey is the verified public key of the SSH CA certificate offered by the client. If no CA certificate
	// was offered this value is nil.
	PubKey(
		username string,
		pubKey string,
		connectionID string,
		remoteAddr net.IP,
		caPubKey *protocol.CACertificate,
	) AuthenticationContext
}

// AuthenticationContext holds the results of an authentication.
type AuthenticationContext interface {
	// Success must return true or false of the authentication was successful / unsuccessful.
	Success() bool
	// Error returns the error that happened during the authentication.
	Error() error
	// Metadata returns a set of metadata entries that have been obtained during the authentication.
	Metadata() map[string]string
}

// NewTestClient creates a new copy of a client usable for testing purposes.
func NewTestClient(cfg config.AuthWebhookClientConfig, logger log.Logger) (Client, error) {
	clientConfig := config.AuthConfig{
		Method:  config.AuthMethodWebhook,
		Webhook: cfg,
	}
	metricsCollector := metrics.New(dummy.New())

	authClient, err := auth.NewHttpAuthClient(
		clientConfig,
		logger,
		metricsCollector,
	)
	if err != nil {
		return nil, err
	}
	return &authClientWrapper{
		authClient,
	}, nil
}

type authClientWrapper struct {
	c auth.Client
}

func (a authClientWrapper) Password(
	username string,
	password []byte,
	connectionID string,
	remoteAddr net.IP,
) AuthenticationContext {
	return a.c.Password(username, password, connectionID, remoteAddr)
}

func (a authClientWrapper) PubKey(
	username string,
	pubKey string,
	connectionID string,
	remoteAddr net.IP,
	caPubKey *protocol.CACertificate,
) AuthenticationContext {
	return a.c.PubKey(username, pubKey, connectionID, remoteAddr, caPubKey)
}
