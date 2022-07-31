package auth

import (
    "context"
    "time"

    "go.containerssh.io/libcontainerssh/metadata"
)

type oidcDeviceFlow struct {
    oidcFlow
    meta metadata.ConnectionAuthPendingMetadata
}

func (o *oidcDeviceFlow) GetAuthorizationURL(ctx context.Context) (
    verificationLink string,
    userCode string,
    expiration time.Duration,
    err error,
) {
    panic("implement me")
}

func (o *oidcDeviceFlow) Verify(ctx context.Context) (string, metadata.ConnectionAuthenticatedMetadata, error) {
    panic("implement me")
}
