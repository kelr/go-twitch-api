package helix

import (
	"context"
	"errors"
	"github.com/google/go-querystring/query"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"golang.org/x/oauth2/twitch"
	"net/http"
	"net/url"
	"io/ioutil"
)

const (
	helixRootURL = "https://api.twitch.tv/helix"
)

// HTTPClient interface for mocking purposes
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Response represents an HTTP response with un-decoded body bytes.
type Response struct {
	Status int
	Header  http.Header
	Data []byte
}

// Client handles communication with the Twitch Helix API.
type Client struct {
	conn         HTTPClient
	ClientID     string
	ClientSecret string
	tokenType    string
}

// NewClient returns a new Twitch Client. If clientID is "", it will not be appended on the request header.
// A client credentials config is established which auto-refreshes OAuth2 access tokens
// Currently ONLY uses Client Credentials flow. Not intended for user access tokens.
// See NewClientUserAuth for user authentication.
func NewClient(clientID string, clientSecret string) (*Client, error) {
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

	return &Client{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		conn:         config.Client(context.Background()),
		tokenType:    "client",
	}, nil
}

// NewClientUserAuth creates a new helix API twitch client with a user token. This token may be obtained with NewUserAuth and TokenExchange, or an existing user token
// may be used instead. The OAuth2 config used to create the token must match. The user token will be automatically refreshed.
func NewClientUserAuth(config *oauth2.Config, userToken *oauth2.Token) (*Client, error) {
	return &Client{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		conn:         config.Client(context.Background(), userToken),
		tokenType:    "user",
	}, nil
}

// Creates a URL with path and the values in params appended onto it
func buildURL(path string, params interface{}) (string, error) {
	targetURL, err := url.Parse(helixRootURL + path)
	if err != nil {
		return "", err
	}

	// Convert optional params to URL queries
	if params != nil {
		qs, err := query.Values(params)
		if err != nil {
			return "", err
		}
		targetURL.RawQuery = qs.Encode()
	}
	return targetURL.String(), nil
}

// Wrapper for a HTTP GET request
func (c *Client) getRequest(path string, params interface{}) (*Response, error) {
	return c.sendRequest(path, params, http.MethodGet)
}

// Wrapper for a HTTP PUT request
func (c *Client) putRequest(path string, params interface{}) (*Response, error) {
	return c.sendRequest(path, params, http.MethodPut)
}

// Wrapper for a HTTP POST request
func (c *Client) postRequest(path string, params interface{}) (*Response, error) {
	return c.sendRequest(path, params, http.MethodPost)
}

// Create and send an HTTP request. Return the decoded JSON value of the HTTP body regardless of status code.
func (c *Client) sendRequest(path string, params interface{}, requestType string) (*Response, error) {
	targetURL, err := buildURL(path, params)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(requestType, targetURL, nil)
	if err != nil {
		return nil, err
	}

	// A client ID is required. The auth token will be added automatically.
	request.Header.Set("Client-ID", c.ClientID)

	// Send the request
	resp, err := c.conn.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	response := new(Response)
	response.Status = resp.StatusCode
	response.Header = resp.Header
	response.Data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return response, nil
}
