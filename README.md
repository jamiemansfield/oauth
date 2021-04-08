oauth
===

oauth is a small client library for interacting with OAuth servers with Go,
licensed under the MIT License.

## Examples

### Refresh Access Token

```go
func main() {
	client := oauth.NewClient(nil,
		"CLIENT_ID_GOES_HERE",
		"CLIENT_SECRET_GOES_HERE")

	token, err := client.RefreshAccessToken(
		"https://login.microsoftonline.com/consumers/oauth2/v2.0/token",
		"REFRESH_TOKEN_GOES_HERE",
		nil,
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("access token: " + token.AccessToken)
}
```

### Device Grant

```go
func main() {
	client := oauth.NewClient(nil,
		"CLIENT_ID_GOES_HERE",
		"CLIENT_SECRET_GOES_HERE")
	deviceGrantService := device.NewService(client)

	auth, err := deviceGrantService.GetAuthorisation(
		"https://login.microsoftonline.com/consumers/oauth2/v2.0/devicecode",
		[]string{"XboxLive.signin", "XboxLive.offline_access"},
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Please visit " + auth.VerificationURI + " and enter the code: " + auth.UserCode)

	token, err := deviceGrantService.PollAccessToken(
		"https://login.microsoftonline.com/consumers/oauth2/v2.0/token",
		auth,
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("access token: " + token.AccessToken)
	fmt.Println("refresh token: " + token.RefreshToken)
}
```

It's worth noting that the client secret isn't needed for the above example, and
won't be sent to the server. It is totally valid to leave the client secret as an
empty string.
