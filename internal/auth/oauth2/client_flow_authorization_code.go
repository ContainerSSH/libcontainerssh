package oauth2

import (
    "context"
    "net/url"
    "strings"
    "time"

    "go.containerssh.io/libcontainerssh/log"
    "go.containerssh.io/libcontainerssh/message"
)

type authorizationCodeFlow struct {
    client         client
    requestScopes  []string
    requiredScopes []string
    state          string
    logger         log.Logger
}

func (a authorizationCodeFlow) GetAuthorizationURL(ctx context.Context) (string, error) {
    endpoint := a.client.authorizationEndpointURL
    l, err := url.Parse(endpoint)
    if err != nil {
        return "", err
    }
    query := l.Query()
    query.Set("response_type", "code")
    query.Set("client_id", a.client.clientID)
    query.Set("scope", strings.Join(a.requestScopes, " "))
    query.Set("state", a.state)
    if a.client.redirectURI != "" {
        query.Set("redirect_uri", a.client.redirectURI)
    }
    l.RawQuery = query.Encode()
    return l.String(), nil
}

func (a authorizationCodeFlow) Verify(
    ctx context.Context,
    state string,
    authorizationCode string,
) (string, []string, error) {
    if a.state != state {
        return "", nil, message.UserMessage(
            message.EAuthOAuth2StateMismatch,
            "Authentication failed, please try again.",
            "The state variable does not match the expected value.",
        )
    }

    var statusCode int
    var lastError error
loop:
    for {
        req := &accessTokenRequest{
            ClientID:     a.client.clientID,
            ClientSecret: a.client.clientSecret,
            Code:         authorizationCode,
        }
        resp := &accessTokenResponse{}
        statusCode, lastError = a.client.httpClient.Post("", req, resp)
        if statusCode != 200 {
            lastError = message.UserMessage(
                message.EAuthOAuth2AccessTokenFetchFailed,
                "Cannot authenticate at this time.",
                "Non-200 status code from oAuth2 access token API (%d; %s; %s).",
                statusCode,
                resp.Error,
                resp.ErrorDescription,
            )
        } else if lastError == nil {
            scopes := strings.Split(resp.Scope, " ")
            return resp.AccessToken, scopes, a.client.checkGrantedScopes(scopes, a.requiredScopes, a.logger)
        }
        a.logger.Debug(lastError)
        select {
        case <-ctx.Done():
            break loop
        case <-time.After(10 * time.Second):
        }
    }
    err := message.WrapUser(
        lastError,
        message.EAuthOAuth2Timeout,
        "Timeout while trying to obtain GitHub authentication data.",
        "Timeout while trying to obtain GitHub authentication data.",
    )
    a.logger.Debug(err)
    return "", nil, err
}
