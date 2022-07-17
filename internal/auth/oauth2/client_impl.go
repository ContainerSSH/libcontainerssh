package oauth2

import (
    "context"
    "encoding/base64"
    "fmt"

    "go.containerssh.io/libcontainerssh/config"
    "go.containerssh.io/libcontainerssh/http"
    "go.containerssh.io/libcontainerssh/log"
    "go.containerssh.io/libcontainerssh/message"
)

func NewClient(
    authorizationEndpointURL string,
    deviceCodeEndpoint string,
    clientID string,
    clientSecret string,
    redirectURI string,
    httpConfig config.HTTPClientConfiguration,
    logger log.Logger,
) (Client, error) {
    cfg := httpConfig
    cfg.RequestEncoding = config.RequestEncodingWWWURLEncoded
    httpClient, err := http.NewClientWithHeaders(
        cfg,
        logger,
        map[string][]string{
            "authorization": {
                "Basic " + base64.StdEncoding.EncodeToString([]byte(clientID+":"+clientSecret)),
            },
        },
        true,
    )
    if err != nil {
        return nil, message.WrapUser(
            err,
            message.EAuthOAuth2HTTPClientCreateFailed,
            "Authentication currently unavailable.",
            "Cannot create authenticator because the token endpoint configuration failed.",
        )
    }

    return &client{
        authorizationEndpointURL: authorizationEndpointURL,
        deviceCodeEndpoint:       deviceCodeEndpoint,
        clientID:                 clientID,
        clientSecret:             clientSecret,
        redirectURI:              redirectURI,
        httpClient:               httpClient,
    }, nil
}

type client struct {
    authorizationEndpointURL string
    deviceCodeEndpoint       string
    clientID                 string
    clientSecret             string
    redirectURI              string
    httpClient               http.Client
}

func (c client) GetAuthorizationCodeFlow(
    ctx context.Context,
    requestScopes []string,
    requiredScopes []string,
    state string,
    logger log.Logger,
) (AuthorizationCodeFlow, bool) {
    return &authorizationCodeFlow{
        client:         c,
        requestScopes:  requestScopes,
        requiredScopes: requiredScopes,
        state:          state,
        logger:         logger,
    }, true
}

func (c client) GetDeviceFlow(
    ctx context.Context,
    requestScopes []string,
    requiredScopes []string,
    state string,
    logger log.Logger,
) (DeviceFlow, bool) {
    return &deviceCodeFlow{
        client:         c,
        requestScopes:  requestScopes,
        requiredScopes: requiredScopes,
        state:          state,
        logger:         logger,
    }, true
}

func (c client) checkGrantedScopes(grantedScopes []string, requiredScopes []string, logger log.Logger) error {
    for _, requiredScope := range requiredScopes {
        scopeGranted := false
        for _, grantedScope := range grantedScopes {
            if grantedScope == requiredScope {
                scopeGranted = true
                break
            }
        }
        if !scopeGranted {
            err := message.UserMessage(
                message.EAuthGitHubRequiredScopeNotGranted,
                fmt.Sprintf("You have not granted us the required %s permission.", requiredScope),
                "The user has not granted the %s permission.",
                requiredScope,
            )
            logger.Debug(err)
            return err
        }
    }
    return nil
}
