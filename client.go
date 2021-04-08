package oauth

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	defaultUserAgent = "oauth/0.1.0 (github.com/jamiemansfield/oauth)"
)

type Client struct {
	httpClient *http.Client

	UserAgent string

	ClientID     string
	ClientSecret string
}

func NewClient(client *http.Client, id string, secret string) *Client {
	if client == nil {
		client = http.DefaultClient
	}
	return &Client{
		httpClient:   client,
		UserAgent:    defaultUserAgent,
		ClientID:     id,
		ClientSecret: secret,
	}
}

func (c *Client) NewRequest(method string, urlStr string, data url.Values) (*http.Request, error) {
	// Create the request
	req, err := http.NewRequest(method, urlStr, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	// Set request headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}

	return req, nil
}

func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = checkResponse(resp)
	if err != nil {
		return resp, err
	}

	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(v)
	}

	return resp, err
}

func checkResponse(r *http.Response) error {
	// 200 is good
	if r.StatusCode == 200 {
		return nil
	}

	var errResp ErrorResponse
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		err := json.Unmarshal(data, &errResp)
		if err != nil {
			return err
		}
	}

	return errResp
}

func (c *Client) postToken(urlStr string, data url.Values) (*AccessTokenResponse, error) {
	req, err := c.NewRequest(http.MethodPost, urlStr, data)
	if err != nil {
		return nil, err
	}

	var authorisation AccessTokenResponse
	_, err = c.Do(req, &authorisation)
	if err != nil {
		return nil, err
	}

	return &authorisation, nil
}

// See https://tools.ietf.org/html/rfc6749#section-6
func (c *Client) RefreshAccessToken(urlStr string, refreshToken string, scope []string) (*AccessTokenResponse, error) {
	data := make(url.Values)
	data.Set("client_id", c.ClientID)
	data.Set("client_secret", c.ClientSecret)
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)
	if scope != nil && len(scope) > 0 {
		data["scope"] = scope
	}

	return c.postToken(urlStr, data)
}
