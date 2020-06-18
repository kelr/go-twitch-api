// Provides easier to use wrapper functions for the Helix API client
package gundyr

import (
	"errors"
	"github.com/kelr/gundyr/helix"
)

type Helix struct {
	client *helix.Client
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
func (c *Helix) IDToUser(userID string) (string, error) {
	opt := &helix.GetUsersOpt{
		ID: []string{userID},
	}

	response, err := c.client.GetUsers(opt)
	if err != nil {
		return "", err
	}

	if len(response.Data) == 0 {
		return "", errors.New("ID: " + userID + " not found")
	}
	return response.Data[0].Login, nil
}

// IdToUser converts multiple user ID strings to multiple username strings.
func (c *Helix) IDsToUser(userIDs []string) ([]string, error) {
	if len(userIDs) > 100 {
		return nil, errors.New("Helix: Cannot request more than 100 user IDs per call.")
	}
	var followers []string

	opt := &helix.GetUsersOpt{
		ID: userIDs,
	}

	response, err := c.client.GetUsers(opt)
	if err != nil {
		return nil, err
	}

	for _, d := range response.Data {
		followers = append(followers, d.Login)
	}
	return followers, nil
}

// UserToId converts a username string to a user ID string.
func (c *Helix) UserToID(username string) (string, error) {
	opt := &helix.GetUsersOpt{
		Login: []string{username},
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

// ViewCount returns the lifetime viewcount of a user by ID.
func (c *Helix) ViewCount(userID string) (int, error) {
	opt := &helix.GetUsersOpt{
		ID: []string{userID},
	}

	response, err := c.client.GetUsers(opt)
	if err != nil {
		return 0, err
	}

	if len(response.Data) == 0 {
		return 0, errors.New("ID: " + userID + " not found")
	}
	return response.Data[0].ViewCount, nil
}

// GetFollowers returns userIDs for all the users following the provided userID.
// "Who is following userID?"
func (c *Helix) GetFollowers(userID string) ([]string, error) {
	var followers []string
	opt := &helix.GetUsersFollowsOpt{
		ToID: userID,
	}
	response, err := c.client.GetUsersFollows(opt)
	if err != nil {
		return followers, err
	}

	for len(response.Data) != 0 {
		for _, d := range response.Data {
			followers = append(followers, d.FromID)
		}

		// Request next page to ensure all the users are found.
		opt = &helix.GetUsersFollowsOpt{
			ToID: userID,
			After: response.Pagination.Cursor,
		}
		response, err = c.client.GetUsersFollows(opt)
		if err != nil {
			return followers, err
		}
	}
	return followers, nil
}
