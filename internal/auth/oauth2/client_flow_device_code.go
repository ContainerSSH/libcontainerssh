package oauth2

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.containerssh.io/libcontainerssh/log"
	"go.containerssh.io/libcontainerssh/message"
)

type deviceCodeFlow struct {
	client         client
	requestScopes  []string
	requiredScopes []string
	state          string
	logger         log.Logger
	deviceCode     string
	interval       time.Duration
}

func (d *deviceCodeFlow) GetAuthorizationURL(ctx context.Context) (verificationLink string, userCode string, expiration time.Duration, err error) {
	req := &deviceRequest{
		ClientID: d.client.clientID,
		Scope:    strings.Join(d.requestScopes, " "),
	}
	var lastError error
	var statusCode int
loop:
	for {
		resp := &deviceResponse{}
		statusCode, lastError = d.client.httpClient.Post(d.client.deviceCodeEndpoint, req, resp)
		if lastError == nil {
			if statusCode == 200 {
				if resp.Interval > 1 {
					d.interval = time.Duration(resp.Interval) * time.Second
				} else {
					d.interval = time.Second
				}
				d.deviceCode = resp.DeviceCode
				return resp.VerificationURI, resp.UserCode, time.Duration(resp.ExpiresIn) * time.Second, nil
			} else {
				switch resp.Error {
				case "slow_down":
					// Let's assume this means that we reached the 50/hr limit. This is currently undocumented.
					lastError = message.UserMessage(
						message.EAuthGitHubDeviceAuthorizationLimit,
						"Cannot authenticate at this time.",
						"GitHub device authorization limit reached (%s).",
						resp.ErrorDescription,
					)
					d.logger.Debug(lastError)
					return "", "", 0, lastError
				}
			}
			lastError = message.UserMessage(
				message.EAuthOAuth2DeviceCodeRequestFailed,
				"Cannot authenticate at this time.",
				"Non-200 status code from oAuth2 device code API (%d; %s; %s).",
				statusCode,
				resp.Error,
				resp.ErrorDescription,
			)
			d.logger.Debug(lastError)
		}
		d.logger.Debug(lastError)
		select {
		case <-time.After(10 * time.Second):
			continue
		case <-ctx.Done():
			break loop
		}
	}
	err = message.WrapUser(
		lastError,
		message.EAuthOAuth2Timeout,
		"Cannot authenticate at this time.",
		"Timeout while trying to obtain a GitHub device code.",
	)
	d.logger.Debug(err)
	return "", "", 0, err
}

func (d deviceCodeFlow) Verify(ctx context.Context) (string, []string, error) {
	var statusCode int
	var lastError error
loop:
	for {
		req := &accessTokenRequest{
			ClientID:   d.client.clientID,
			DeviceCode: d.deviceCode,
			GrantType:  "urn:ietf:params:oauth:grant-type:device_code",
		}
		resp := &accessTokenResponse{}
		statusCode, lastError = d.client.httpClient.Post("", req, resp)
		if statusCode != 200 {
			if resp.Error == "authorization_pending" {
				lastError = message.NewMessage(
					message.EAuthOAuth2AuthorizationPending,
					"User authorization still pending, retrying in %d seconds.",
					d.interval,
				)
			} else {
				lastError = message.UserMessage(
					message.EAuthGitHubAccessTokenFetchFailed,
					"Cannot authenticate at this time.",
					"Non-200 status code from GitHub access token API (%d; %s; %s).",
					statusCode,
					resp.Error,
					resp.ErrorDescription,
				)
			}
		} else if lastError == nil {
			switch resp.Error {
			case "authorization_pending":
				lastError = message.UserMessage(message.EAuthOAuth2AuthorizationPending, "Authentication is still pending.", "The user hasn't completed the authentication process.")
			case "slow_down":
				if resp.Interval > 30 {
					// Assume we have exceeded the hourly rate limit, let's fall back to authorization code flow.
					return "", nil, message.UserMessage(
						message.EAuthDeviceFlowRateLimitExceeded,
						"Cannot authenticate at this time. Please try again later.",
						"Rate limit for device flow exceeded, attempting authorization code flow.",
					)
				}
			case "expired_token":
				return "", nil, fmt.Errorf("BUG: expired token during device flow authentication")
			case "unsupported_grant_type":
				return "", nil, message.UserMessage(
					message.EAuthOAuth2UnsupportedGrantType,
					"Cannot authenticate at this time. Please try again later.",
					"The oAuth2 server responded with an 'unsupported grant type'. Check your oAuth2 configuration and enable only the authentication types your oAuth2 server supports.",
				)
			case "incorrect_client_credentials":
				// User entered the incorrect device code
				return "", nil, message.UserMessage(
					message.EAuthIncorrectClientCredentials,
					"Authentication failed",
					"User entered incorrect device code.",
				)
			case "incorrect_device_code":
				// User entered the incorrect device code
				return "", nil, message.UserMessage(
					message.EAuthFailed,
					"Authentication failed",
					"User entered incorrect device code",
				)
			case "access_denied":
				// User hit don't authorize
				return "", nil, message.UserMessage(
					message.EAuthFailed,
					"Authentication failed",
					"User canceled oAuth2 authentication",
				)
			case "":
				scopes := strings.Split(resp.Scope, " ")
				return resp.AccessToken, scopes, d.client.checkGrantedScopes(scopes, d.requiredScopes, d.logger)
			}
		}
		d.logger.Debug(lastError)
		interval := d.interval
		if resp.Interval != 0 {
			interval = time.Duration(int(resp.Interval)) * time.Second
		}
		select {
		case <-ctx.Done():
			break loop
		case <-time.After(interval):
		}
	}
	err := message.WrapUser(
		lastError,
		message.EAuthOAuth2Timeout,
		"Timeout while trying to obtain GitHub authentication data.",
		"Timeout while trying to obtain GitHub authentication data.",
	)
	d.logger.Debug(err)
	return "", nil, err
}
