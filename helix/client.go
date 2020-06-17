package helix

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/go-querystring/query"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"golang.org/x/oauth2/twitch"
	"io"
	"net/http"
	"net/url"
)

const (
	helixRootURL = "https://api.twitch.tv/helix"
)

// HTTPClient interface for mocking purposes
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// TwitchClient handles communication with the Twitch Helix API.
type TwitchClient struct {
	conn         HTTPClient
	ClientID     string
	ClientSecret string
	tokenType    string
}

// NewTwitchClient returns a new Twitch Client. If clientID is "", it will not be appended on the request header.
// A client credentials config is established which auto-refreshes OAuth2 access tokens
// Currently ONLY uses Client Credentials flow. Not intended for user access tokens.
// See NewTwitchClientUserAuth for user authentication.
func NewTwitchClient(clientID string, clientSecret string) (*TwitchClient, error) {
	if clientID == "" {
		return nil, errors.New("A Client ID must be provided to create a twitch client")
	}
	if clientSecret == "" {
		return nil, errors.New("A Client secret must be provided to create a twitch client")
	}
	config := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     twitch.Endpoint.TokenURL,
	}

	_, err := config.Token(context.Background())
	if err != nil {
		return nil, err
	}

	return &TwitchClient{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		conn:         config.Client(context.Background()),
		tokenType:    "client",
	}, nil
}

// NewTwitchClientUserAuth creates a new helix API twitch client with a user token. This token may be obtained with NewUserAuth and TokenExchange, or an existing user token
// may be used instead. The OAuth2 config used to create the token must match. The user token will be automatically refreshed.
func NewTwitchClientUserAuth(config *oauth2.Config, userToken *oauth2.Token) (*TwitchClient, error) {
	return &TwitchClient{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		conn:         config.Client(context.Background(), userToken),
		tokenType:    "user",
	}, nil
}

// Creates a URL with path and the values in params appended onto it
func buildURL(path string, params interface{}) (*url.URL, error) {
	targetURL, err := url.Parse(helixRootURL + path)
	if err != nil {
		return nil, err
	}

	// Convert optional params to URL queries
	if params != nil {
		qs, err := query.Values(params)
		if err != nil {
			return nil, err
		}
		targetURL.RawQuery = qs.Encode()
	}
	return targetURL, nil
}

// Create and send an HTTP request. Return the decoded JSON value of the HTTP body regardless of status code.
func (client *TwitchClient) sendRequest(path string, params interface{}, result interface{}, requestType string) (*http.Response, error) {
	targetURL, err := buildURL(path, params)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(requestType, targetURL.String(), nil)
	if err != nil {
		return nil, err
	}

	// A client ID is required. The auth token will be added automatically.
	request.Header.Set("Client-ID", client.ClientID)

	// Send the request
	resp, err := client.conn.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Decode the response
	err = json.NewDecoder(resp.Body).Decode(result)
	if err == io.EOF {
		err = nil
	}

	// TODO: Check response code
	return resp, nil
}
