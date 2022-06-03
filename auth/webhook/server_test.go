package webhook_test

import (
	"context"
	"os"
	"time"

    auth2 "go.containerssh.io/libcontainerssh/auth"
    "go.containerssh.io/libcontainerssh/auth/webhook"
    "go.containerssh.io/libcontainerssh/config"
    "go.containerssh.io/libcontainerssh/log"
    "go.containerssh.io/libcontainerssh/metadata"
    "go.containerssh.io/libcontainerssh/service"
)

// myAuthReqHandler is your handler for authentication requests.
type myAuthReqHandler struct {
}

// OnPassword will be called when the user requests password authentication.
func (m *myAuthReqHandler) OnPassword(
	meta metadata.ConnectionAuthPendingMetadata,
	password []byte,
) (
	success bool,
	metadata metadata.ConnectionAuthenticatedMetadata,
	err error,
) {
	return true, meta.Authenticated(meta.Username), nil
}

// OnPubKey will be called when the user requests public key authentication.
func (m *myAuthReqHandler) OnPubKey(
	meta metadata.ConnectionAuthPendingMetadata,
	publicKey auth2.PublicKey,
) (
	success bool,
	metadata metadata.ConnectionAuthenticatedMetadata,
	err error,
) {
	return true, meta.Authenticated(meta.Username), nil
}

// OnAuthorization will be called after login in non-webhook auth handlers to verify the user is authorized to login
func (m *myAuthReqHandler) OnAuthorization(
	meta metadata.ConnectionAuthenticatedMetadata,
) (
	success bool,
	metadata metadata.ConnectionAuthenticatedMetadata,
	err error,
) {
	return true, meta, nil
}

// ExampleNewServer demonstrates how to set up an authentication server.
func ExampleNewServer() {
	// Set up a logger.
	logger := log.MustNewLogger(config.LogConfig{
		Level:       config.LogLevelWarning,
		Format:      config.LogFormatText,
		Destination: config.LogDestinationStdout,
		Stdout:      os.Stdout,
	})

	// Create a new auth webhook server.
	srv, err := webhook.NewServer(
		config.HTTPServerConfiguration{
			Listen: "0.0.0.0:8001",
		},
		// Pass your handler here.
		&myAuthReqHandler{},
		logger,
	)
	if err != nil {
		// Handle error
		panic(err)
	}

	// Set up and run the web server service.
	lifecycle := service.NewLifecycle(srv)

	go func() {
		//Ignore error, handled later.
		_ = lifecycle.Run()
	}()

	// Sleep for 30 seconds as a test.
	time.Sleep(30 * time.Second)

	// Set up a shutdown context to give a deadline for graceful shutdown.
	shutdownContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	// Stop the server.
	lifecycle.Stop(shutdownContext)

	// Wait for the server to stop.
	lastError := lifecycle.Wait()
	if lastError != nil {
		// Server stopped abnormally.
		panic(lastError)
	}

	// Output:
}
