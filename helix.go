// Provides easier to use wrapper functions for the Helix API client
package gundyr

import (
	"errors"
	"github.com/kelr/gundyr/helix"
)

// Interface to allow for mocking a Helix Client.
type helixClient interface {
	GetUsers(opt *helix.GetUsersOpt) (*helix.GetUsersResponse, error)
	GetUsersFollows(opt *helix.GetUsersFollowsOpt) (*helix.GetUsersFollowsResponse, error)
}

// Helix is a wrapper over a HelixClient.
type Helix struct {
	client helixClient
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

// IDToUser converts a user ID string to a username string.
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

// IDsToUser converts multiple user ID strings to multiple username strings.
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
		return followers, err
	}

	for _, d := range response.Data {
		followers = append(followers, d.Login)
	}
	return followers, nil
}

// UserToID converts a username string to a user ID string.
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

	// Drain all the followers by checking each page until there are none left.
	for len(response.Data) > 0 {
		for _, d := range response.Data {
			followers = append(followers, d.FromID)
		}

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
