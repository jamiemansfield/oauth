package oauth

// Possible error types as defined by RFC 6749 (OAuth 2.0).
//
// See https://tools.ietf.org/html/rfc6749#section-5.2
const (
	ErrorInvalidRequest       = "invalid_request"
	ErrorInvalidClient        = "invalid_client"
	ErrorInvalidGrant         = "invalid_grant"
	ErrorUnauthorisedClient   = "unauthorized_client"
	ErrorUnsupportedGrantType = "unsupported_grant_type"
	ErrorInvalidScope         = "invalid_scope"
)

// ErrorResponse is a model of the 'Error Response', as defined by the RFC
// 6749 (OAuth 2.0) specification.
//
// See https://tools.ietf.org/html/rfc6749#section-5.2
type ErrorResponse struct {
	Code        string `json:"error"`
	Description string `json:"error_description"`
	URI         string `json:"error_uri"`
}

var _ error = (*ErrorResponse)(nil)

func (r ErrorResponse) Error() string {
	var message string
	message = r.Code

	if r.Description != "" {
		message += ": " + r.Description
	}

	if r.URI != "" {
		message += " (" + r.URI + ")"
	}

	return message
}
