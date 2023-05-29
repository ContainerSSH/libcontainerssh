package oauth2

type accessTokenRequest struct {
	GrantType    string `json:"grant_type,omitempty" schema:"grant_type"`
	Code         string `json:"code,omitempty" schema:"code"`
	DeviceCode   string `json:"device_code,omitempty" schema:"device_code"`
	ClientID     string `json:"client_id" schema:"client_id,required"`
	ClientSecret string `json:"client_secret" schema:"client_secret"`
	RedirectURI  string `json:"redirect_uri" schema:"redirect_uri"`
}

type accessTokenResponse struct {
	AccessToken      string `json:"access_token,omitempty"`
	TokenType        string `json:"token_type,omitempty"`
	ExpiresIn        int    `json:"expires_in,omitempty"`
	RefreshToken     string `json:"refresh_token,omitempty"`
	Scope            string `json:"scope,omitempty"`
	Interval         uint   `json:"interval,omitempty" yaml:"interval"`
	Error            string `json:"error,omitempty"`
	ErrorDescription string `json:"error_description,omitempty"`
	ErrorURI         string `json:"error_uri,omitempty"`
}

type deviceRequest struct {
	ClientID string `schema:"client_id"`
	Scope    string `schema:"scope"`
}

type deviceResponse struct {
	DeviceCode       string `json:"device_code"`
	UserCode         string `json:"user_code"`
	VerificationURI  string `json:"verification_uri"`
	ExpiresIn        uint   `json:"expires_in" yaml:"expires_in"`
	Interval         uint   `json:"interval" yaml:"interval"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	ErrorURI         string `json:"error_uri"`
}
