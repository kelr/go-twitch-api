// Package gundyr provides an interface to access the Helix Twitch API and Twich PubSub.
package gundyr

import (
	"errors"
	"github.com/kelr/gundyr/helix"
	"golang.org/x/oauth2"
)

var pageMem = make(map[string]int)

// Interface to allow for mocking a Helix Client.
type helixClient interface {
	GetUsers(opt *helix.GetUsersOpt) (*helix.GetUsersResponse, error)
	GetUsersFollows(opt *helix.GetUsersFollowsOpt) (*helix.GetUsersFollowsResponse, error)
	GetClips(opt *helix.GetClipsOpt) (*helix.GetClipsResponse, error)
}

// HelixConfig represents configuration options available to a Client.
type HelixConfig struct {
	ClientID     string
	ClientSecret string
	Scopes       []string
	RedirectURI  string
	Token        *oauth2.Token
}

// Helix is a wrapper over a HelixClient. See https://godoc.org/github.com/kelr/gundyr/helix for the underlying HelixClient.
type Helix struct {
	client helixClient
}

// NewHelix returns returns a client credentials Helix API client wrapper
func NewHelix(cfg *HelixConfig) (*Helix, error) {
	c := helix.Config(*cfg)
	client, err := helix.NewClient(&c)
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
// Will accept a maximum of 100 IDs. Requests for more than 100 IDs should call this function
// in chunks.
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

// GetUserEmail returns the e-mail address of the user by username.
// The user access token must have scope user:read:email for the username provided.
func (c *Helix) GetUserEmail(username string) (string, error) {
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

	return response.Data[0].Email, nil
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
			ToID:  userID,
			After: response.Pagination.Cursor,
		}

		response, err = c.client.GetUsersFollows(opt)
		if err != nil {
			return followers, err
		}
	}
	return followers, nil
}

func (c *Helix) GetAllClips(broadcasterID string, after string) ([]helix.GetClipsData, error) {
	var clips []helix.GetClipsData

	opt := &helix.GetClipsOpt{
		BroadcasterID: broadcasterID,
		StartedAt:     after,
		First:         100,
	}

	response, err := c.client.GetClips(opt)
	if err != nil {
		return nil, err
	}

	// Drain all the clips by checking each page until there are none left.
	for len(response.Data) > 0 {
		if _, ok := pageMem[response.Pagination.Cursor]; ok {
			pageMem = make(map[string]int)
			return clips, nil
		}

		for _, c := range response.Data {
			clips = append(clips, c)
		}

		pageMem[response.Pagination.Cursor] = 1

		opt = &helix.GetClipsOpt{
			BroadcasterID: broadcasterID,
			After:         response.Pagination.Cursor,
			StartedAt:     after,
			First:         100,
		}

		response, err = c.client.GetClips(opt)
		if err != nil {
			return nil, err
		}
	}
	pageMem = make(map[string]int)
	return clips, nil
}
