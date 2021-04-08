package oauth

// AccessTokenResponse is a model of the 'Successful Response' for
// 'Issuing an Access Token'.
//
// See https://tools.ietf.org/html/rfc6749#section-5.1
type AccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}
