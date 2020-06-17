// Provides easier to use wrapper functions for the Helix API client
package gundyr

import (
	"errors"
	"github.com/kelr/gundyr/helix"
)

type Helix struct {
	client *helix.HelixClient
}

// NewHelix returns returns a client credentials Helix API client wrapper
func NewHelix(clientID string, clientSecret string) (*Helix, error) {
	client, err := helix.NewClient(clientID, clientSecret)
	if err != nil {
		return nil, err
	}
	return &Helix{
		client,
	}, nil
}

// IdToUser converts a user ID string to a username string.
func (c *Helix) IdToUser(userId string) (string, error) {
	opt := &helix.GetUsersOpt{
		ID: userId,
	}

	response, err := c.client.GetUsers(opt)
	if err != nil {
		return "", err
	}

	if len(response.Data) == 0 {
		return "", errors.New("ID: " + userId + " not found")
	}
	return response.Data[0].Login, nil
}

// UserToId converts a username string to a user ID string.
func (c *Helix) UserToId(username string) (string, error) {
	opt := &helix.GetUsersOpt{
		Login: username,
	}

	response, err := c.client.GetUsers(opt)
	if err != nil {
		return "", err
	}

	if len(response.Data) == 0 {
		return "", errors.New("User: " + username + " not found")
	}

	return response.Data[0].ID, nil
}
