// Package helix provides a HTTP client to communicate with the Twitch Helix API endpoints.
package helix

import (
	"context"
	"errors"
	"github.com/google/go-querystring/query"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"golang.org/x/oauth2/twitch"
	"io/ioutil"
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

// Response represents an HTTP response with un-decoded body bytes.
type Response struct {
	Status int
	Header http.Header
	Data   []byte
}

// Client handles communication with the Twitch Helix API.
type Client struct {
	conn      HTTPClient
	config    *Config
	tokenType string
}

// Config represents configuration options available to a Client.
type Config struct {
	ClientID     string
	ClientSecret string
	Scopes       []string
	RedirectURI  string
	Token        *oauth2.Token
}

// NewClient returns a new Helix Client depending on options provided by cfg.
// If the Token field of the config is nil, the client will attempt create an app access token
// using the 2-legged OAuth2 client credentials flow.
// If a Token is provided, the client will attempt to use the token as a user access token.
func NewClient(cfg *Config) (*Client, error) {
	if cfg.ClientID == "" {
		return nil, errors.New("A Client ID must be provided to create a twitch client")
	}

	if cfg.ClientSecret == "" {
		return nil, errors.New("A Client secret must be provided to create a twitch client")
	}

	c := new(Client)
	c.config = cfg

	var err error
	if cfg.Token == nil {
		c.conn, err = newAppAccessClient(cfg)
		c.tokenType = "app"
		if err != nil {
			return nil, err
		}
	} else {
		c.conn = newUserAccessClient(cfg)
		c.tokenType = "user"
	}
	return c, nil
}

func newAppAccessClient(cfg *Config) (*http.Client, error) {
	c := &clientcredentials.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		TokenURL:     twitch.Endpoint.TokenURL,
	}

	_, err := c.Token(context.Background())
	if err != nil {
		return nil, err
	}

	return c.Client(context.Background()), nil
}

func newUserAccessClient(cfg *Config) *http.Client {
	c := &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		Scopes:       cfg.Scopes,
		Endpoint:     twitch.Endpoint,
		RedirectURL:  cfg.RedirectURI,
	}
	return c.Client(context.Background(), cfg.Token)
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

func (c *Client) hasScope(scope string) bool {
	for _, s := range c.config.Scopes {
		if s == scope {
			return true
		}
	}
	return false
}

// Wrapper for a HTTP GET request
func (c *Client) getRequest(path string, params interface{}) (*Response, error) {
	request, err := c.buildRequest(path, params, http.MethodGet)
	if err != nil {
		return nil, err
	}
	return c.sendRequest(request)
}

// Wrapper for a HTTP PUT request
func (c *Client) putRequest(path string, params interface{}) (*Response, error) {
	request, err := c.buildRequest(path, params, http.MethodPut)
	if err != nil {
		return nil, err
	}
	return c.sendRequest(request)
}

// Wrapper for a HTTP POST request
func (c *Client) postRequest(path string, params interface{}) (*Response, error) {
	request, err := c.buildRequest(path, params, http.MethodPost)
	if err != nil {
		return nil, err
	}
	return c.sendRequest(request)
}

// Create an HTTP request.
func (c *Client) buildRequest(path string, params interface{}, requestType string) (*http.Request, error) {
	targetURL, err := buildURL(path, params)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(requestType, targetURL, nil)
	if err != nil {
		return nil, err
	}

	// A client ID is required. The auth token will be added automatically.
	request.Header.Set("Client-ID", c.config.ClientID)
	return request, nil
}

// Create and send an HTTP request. Return the decoded JSON value of the HTTP body regardless of status code.
func (c *Client) sendRequest(request *http.Request) (*Response, error) {
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
