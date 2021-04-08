package device

// AuthResponse is a model of the 'Device Authorisation Response'.
//
// See https://tools.ietf.org/html/rfc8628#section-3.2
type AuthResponse struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURI string `json:"verification_uri"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
}

// Possible error types that can occur when requesting an access token
// using a device grant.
//
// See https://tools.ietf.org/html/rfc8628#section-3.5
const (
	ErrorAuthorisationPending = "authorization_pending"
	ErrorSlowDown             = "slow_down"
	ErrorAccessDenied         = "access_denied"
	ErrorExpiredToken         = "expired_token"
)
