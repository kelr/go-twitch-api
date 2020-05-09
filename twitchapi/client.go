package twitchapi

import (
	"context"
	"encoding/json"
	"fmt"
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

// Handles communication with the Twitch API.
type TwitchClient struct {
	conn         *http.Client
	ClientID     string
	ClientSecret string
	tokenType    string
}

// Returns a new Twitch Client. If clientID is "", it will not be appended on the request header.
// A client credentials config is established which auto-refreshes OAuth2 access tokens
// Currently ONLY uses Client Credentials flow. Not intended for user access tokens.
func NewTwitchClient(clientID string, clientSecret string) (*TwitchClient, error) {
	config := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     twitch.Endpoint.TokenURL,
	}

	_, err := config.Token(context.Background())
	if err != nil {
		fmt.Println("Error in getting a token:", err)
		return nil, err
	}

	return &TwitchClient{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		conn:         config.Client(context.Background()),
		tokenType:    "client",
	}, nil
}

func NewUserAuth(clientID string, clientSecret string, redirectURI string, scopes *[]string) (*oauth2.Config, string) {
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       *scopes,
		Endpoint:     twitch.Endpoint,
		RedirectURL:  redirectURI,
	}
	return config, config.AuthCodeURL("state", oauth2.AccessTypeOffline)
}

func TokenExchange(config *oauth2.Config, authCode string) (*oauth2.Token, error) { 
	token, err := config.Exchange(context.Background(), authCode)
	if err != nil {
		fmt.Println("Error in obtaining user token:", err)
		return nil, err
	}
	return token, nil
}

func NewTwitchClientUserAuth(config *oauth2.Config, userToken *oauth2.Token) (*TwitchClient, error) {
	return &TwitchClient{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		conn:         config.Client(context.Background(), userToken),
		tokenType:    "user",
	}, nil
}

// Create and send an HTTP request.
func (client *TwitchClient) sendRequest(path string, params interface{}, result interface{}, requestType string) (*http.Response, error) {
	targetUrl, err := url.Parse(helixRootURL + path)
	if err != nil {
		return nil, err
	}

	// Convert optional params to URL queries
	if params != nil {
		qs, err := query.Values(params)
		if err != nil {
			return nil, err
		}
		targetUrl.RawQuery = qs.Encode()
	}

	request, err := http.NewRequest(requestType, targetUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	// Add on optional headers
	if client.ClientID != "" {
		request.Header.Set("Client-ID", client.ClientID)
	}

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
