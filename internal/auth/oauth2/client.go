package oauth2

import (
    "context"
    "time"

    "go.containerssh.io/libcontainerssh/log"
)

// Client is a generic oAuth2 client implementation.
type Client interface {
    // GetAuthorizationCodeFlow returns a flow that takes the user to a website, and after login the user is returned
    // with an authorization code. The authorization code is used to get a permanent access token.
    GetAuthorizationCodeFlow(
        ctx context.Context,
        requestScopes []string,
        requiredScopes []string,
        state string,
        logger log.Logger,
    ) (AuthorizationCodeFlow, bool)

    // GetDeviceFlow returns a flow that takes a user to a website where the user needs to enter a code.
    // In the meantime, ContainerSSH polls the oAuth2 server and continues the login when the user has entered
    // the code.
    GetDeviceFlow(
        ctx context.Context,
        requestScopes []string,
        requiredScopes []string,
        state string,
        logger log.Logger,
    ) (DeviceFlow, bool)
}

// AuthorizationCodeFlow is an oAuth23 flow where the user is given a link to click, and on return, must submit an
// authorization code. On traditional applications the submission of the code is done automatically, but with SSH that
// is not an option because many SSH clients don't wait for the process to be complete.
type AuthorizationCodeFlow interface {
    // GetAuthorizationURL returns the authorization URL a user should be redirected to begin the login process.
    //
    // Parameters:
    //
    // - A context that can be used to cancel fetching the process of fetching the URL
    //
    // Returns:
    //
    // - A link the user should be redirected to for login
    // - An error if fetching the authorization URL failed
    GetAuthorizationURL(ctx context.Context) (string, error)

    // Verify verifies the authorizationCode with the OAuth2 server and obtains an access token. It returns the access
    // token if successful, or an error otherwise.
    //
    // Parameters:
    //
    // - A context that can be used to cancel the verification
    // - A state string used to sync the authentication
    // - The authorization code the user received after returning from the oAuth2 flow
    //
    // Returns:
    //
    // - The access token obtained from the authorization code
    // - The list of granted scopes
    // - An error if the process failed
    Verify(ctx context.Context, state string, authorizationCode string) (
        string,
        []string,
        error,
    )
}

// DeviceFlow is an oAuth2 flow where the user is given a link to click together with a device code they need to enter
// on the website. Once the code is entered, the authentication proceeds.
type DeviceFlow interface {
    // GetAuthorizationURL returns the authorization URL a user should be redirected to begin the login process.
    //
    // Parameters:
    //
    // - A context that can be used to cancel fetching the process of fetching the URL
    //
    // Returns:
    //
    // - A link the user should be redirected to for login
    // - A device code the user should enter
    // - An expiration time after which the device code is invalid
    // - An error if the authorization URL cannot be fetched.
    GetAuthorizationURL(
        ctx context.Context,
    ) (
        verificationLink string,
        userCode string,
        expiration time.Duration,
        err error,
    )

    // Verify starts polling the OAuth2 server if the authorization has been completed. The ctx parameter contains a
    // context that can be used to stop the polling.
    //
    // Parameters:
    //
    // - A context that can be used to cancel the verification process
    //
    // Returns:
    //
    // - An access token
    // - A list of granted scopes
    // - The error that happened if fetching the access token failed
    Verify(
        ctx context.Context,
    ) (
        string,
        []string,
        error,
    )
}
