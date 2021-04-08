package device

import (
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/jamiemansfield/oauth"
)

type Service struct {
	client *oauth.Client
}

func NewService(client *oauth.Client) *Service {
	return &Service{client: client}
}

// See https://tools.ietf.org/html/rfc8628#section-3.2
func (s *Service) GetAuthorisation(urlStr string, scope []string) (*AuthResponse, error) {
	data := make(url.Values)
	data.Set("client_id", s.client.ClientID)
	if scope != nil && len(scope) > 0 {
		data["scope"] = scope
	}

	req, err := s.client.NewRequest(http.MethodPost, urlStr, data)
	if err != nil {
		return nil, err
	}

	var authorisation AuthResponse
	_, err = s.client.Do(req, &authorisation)
	if err != nil {
		return nil, err
	}

	return &authorisation, nil
}

// See https://tools.ietf.org/html/rfc8628#section-3.4
func (s *Service) getAccessToken(urlStr string, deviceCode string) (*oauth.AccessTokenResponse, error) {
	data := make(url.Values)
	data.Set("grant_type", "urn:ietf:params:oauth:grant-type:device_code")
	data.Set("device_code", deviceCode)
	data.Set("client_id", s.client.ClientID)

	req, err := s.client.NewRequest(http.MethodPost, urlStr, data)
	if err != nil {
		return nil, err
	}

	var response oauth.AccessTokenResponse
	_, err = s.client.Do(req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// See https://tools.ietf.org/html/rfc8628#section-3.5
func (s *Service) PollAccessToken(urlStr string, authResponse *AuthResponse) (*oauth.AccessTokenResponse, error) {
	// Default to 5 second interval as per the specification
	intervalRaw := authResponse.Interval
	if intervalRaw == 0 {
		intervalRaw = 5
	}

	now := time.Now()
	interval := time.Duration(intervalRaw) * time.Second
	expires := now.Add(time.Duration(authResponse.ExpiresIn) * time.Second)

	for {
		time.Sleep(interval)
		if time.Now().After(expires) {
			return nil, errors.New("oauth/device: Timeout occurred while polling for access token")
		}

		resp, err := s.getAccessToken(urlStr, authResponse.DeviceCode)
		if err == nil {
			return resp, nil
		}

		// Handle errors
		if err, ok := err.(oauth.ErrorResponse); ok {
			if err.Code == ErrorSlowDown {
				// Increase interval by 5 seconds
				interval += time.Duration(5) * time.Second
			} else if err.Code != ErrorAuthorisationPending {
				return nil, err
			}
		} else {
			// todo: should we return unknown errors?
			return nil, err
		}
	}
}
