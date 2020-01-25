package twitchapi

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

const (
	helixRootURL = "https://api.twitch.tv/helix"
)

// Handles communication with the Twitch API.
type TwitchClient struct {
	conn     *http.Client
	// Twitch Client ID
	ClientID string
}

// Returns a new Twitch Client. If clientID is "", it will not be appended on the request header.
func NewTwitchClient(clientID string) *TwitchClient {
	return &TwitchClient{
		ClientID: clientID,
		conn:     http.DefaultClient,
	}
}

// Create and send an HTTP request. 
func (client *TwitchClient) sendRequest(path string, params interface{}, result interface{}) (*http.Response, error) {
	targetUrl, err := url.Parse(helixRootURL + path)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("GET", targetUrl.String(), nil)
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

	err = json.NewDecoder(resp.Body).Decode(result)
	if err == io.EOF {
		err = nil
	}

	// TODO: Check response code

	return resp, nil
}
