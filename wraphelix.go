// Provides easier to use wrapper functions for the Helix API client
package twitchapi

import (
	"errors"
	"github.com/kelr/go-twitch-api/helix"
)

type HelixClient struct {
	*helix.TwitchClient
}

// Creates a client credentials helix API client wrapper
func NewHelixClient(clientID string, clientSecret string) (*HelixClient, error) {
	client, err := helix.NewTwitchClient(clientID, clientSecret)
	if err != nil {
		return nil, err
	}
	return &HelixClient{
		client,
	}, nil
}

// Converts a user ID string to a username string
func (c *HelixClient) IdToUser(userId string) (string, error) {
	opt := &helix.GetUsersOpt{
		ID: userId,
	}

	response, err := c.GetUsers(opt)
	if err != nil {
		return "", err
	}

	if len(response.Data) == 0 {
		return "", errors.New("ID: " + userId + " not found")
	}
	return response.Data[0].Login, nil
}

// Converts a username string to a user ID string
func (c *HelixClient) UserToId(username string) (string, error) {
	opt := &helix.GetUsersOpt{
		Login: username,
	}

	response, err := c.GetUsers(opt)
	if err != nil {
		return "", err
	}

	if len(response.Data) == 0 {
		return "", errors.New("User: " + username + " not found")
	}

	return response.Data[0].ID, nil
}
