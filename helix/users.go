package helix

import (
	"time"
)

const (
	getUsersPath                 = "/users"
	getUsersFollowsPath          = "/users/follows"
	getUsersExtensionsPath       = "/users/extensions/list"
	getUsersActiveExtensionsPath = "/users/extensions"
)

// Defines the options available for Get Users
type GetUsersOpt struct {
	ID    string `url:"id,omitempty"`
	Login string `url:"login,omitempty"`
}

// Response structure for a Get Users command
type GetUsersResponse struct {
	Data []struct {
		ID              string `json:"id,omitempty"`
		Login           string `json:"login,omitempty"`
		DisplayName     string `json:"display_name,omitempty"`
		Type            string `json:"type,omitempty"`
		BroadcasterType string `json:"broadcaster_type,omitempty"`
		Description     string `json:"description,omitempty"`
		ProfileImageUrl string `json:"profile_image_url,omitempty"`
		OfflineImageUrl string `json:"offline_image_url,omitempty"`
		ViewCount       string `json:"view_count,omitempty"`
		Email           string `json:"email,omitempty"`
	} `json:"data,omitempty"`
}

// Return a slice representing the information for the requested user(s)
//
// https://dev.twitch.tv/docs/api/reference#get-users
func (client *TwitchClient) GetUsers(opt *GetUsersOpt) (*GetUsersResponse, error) {
	data := new(GetUsersResponse)
	_, err := client.sendRequest(getUsersPath, opt, data, "GET")
	if err != nil {
		return nil, err
	}
	return data, err
}

// Defines the options available for Get Users Follows
type GetUsersFollowsOpt struct {
	After  string `url:"after,omitempty"`
	First  int    `url:"first,omitempty"`
	FromID string `url:"from_id,omitempty"`
	ToID   string `url:"to_id,omitempty"`
}

// Response structure for a Get Users Follows command
type GetUsersFollowsResponse struct {
	Total int `json:"total,omitempty"`
	Data  []struct {
		FollowedAt time.Time `json:"followed_at,omitempty"`
		FromID     string    `json:"from_id,omitempty"`
		FromName   string    `json:"from_name,omitempty"`
		ToID       string    `json:"to_id,omitempty"`
		ToName     string    `json:"to_name,omitempty"`
	} `json:"data,omitempty"`
	Pagination struct {
		Cursor string `json:"cursor,omitempty"`
	} `json:"pagination,omitempty"`
}

// Return a slice representing the followers from ids or followers to ids
//
// https://dev.twitch.tv/docs/api/reference#get-users-follows
func (client *TwitchClient) GetUsersFollows(opt *GetUsersFollowsOpt) (*GetUsersFollowsResponse, error) {
	data := new(GetUsersFollowsResponse)
	_, err := client.sendRequest(getUsersFollowsPath, opt, data, "GET")
	if err != nil {
		return nil, err
	}
	return data, err
}

// Defines the options available for Update User
type UpdateUserOpt struct {
	Description string `url:"description"`
}

// Updates the description of a user. Requires a user token for the user to be updated.
// Requires scope: user:edit
//
// https://dev.twitch.tv/docs/api/reference#update-user
func (client *TwitchClient) UpdateUser(opt *UpdateUserOpt) (*GetUsersResponse, error) {
	data := new(GetUsersResponse)
	_, err := client.sendRequest(getUsersPath, opt, data, "PUT")
	if err != nil {
		return nil, err
	}
	return data, err
}

// Response structure for a Get Users Extensions command
type GetUserExtensionsResponse struct {
	Data []struct {
		ID          string   `json:"id,omitempty"`
		Version     string   `json:"version,omitempty"`
		Name        string   `json:"name,omitempty"`
		CanActivate bool     `json:"can_activate,omitempty"`
		Type        []string `json:"type,omitempty"`
	} `json:"data,omitempty"`
}

// Returns a list of active and inactive extensions for a user identified by the user token
// Requires scope user:read:broadcast
//
// https://dev.twitch.tv/docs/api/reference#get-users-follows
func (client *TwitchClient) GetUserExtensions() (*GetUserExtensionsResponse, error) {
	data := new(GetUserExtensionsResponse)
	_, err := client.sendRequest(getUsersExtensionsPath, nil, data, "GET")
	if err != nil {
		return nil, err
	}
	return data, err
}

type GetUserActiveExtensionsOpt struct {
	UserID string `url:"user_id"`
}

type ActiveExtension struct {
	Active  bool   `json:"active"`
	ID      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
	X       int    `json:"x,omitempty"`
	Y       int    `json:"y,omitempty"`
}

// Response structure for a Get Users Extensions command
type GetUserActiveExtensionsResponse struct {
	Data struct {
		Component map[string]ActiveExtension `json:"component,omitempty"`
		Overlay   map[string]ActiveExtension `json:"overlay,omitempty"`
		Panel     map[string]ActiveExtension `json:"panel,omitempty"`
	} `json:"data,omitempty"`
}

// Returns a list of active and inactive extensions for a user identified by the user token
// Requires scope user:read:broadcast or user:edit:broadcast
//
// https://dev.twitch.tv/docs/api/reference#get-user-active-extensions
func (client *TwitchClient) GetUserActiveExtensions(opt *GetUserActiveExtensionsOpt) (*GetUserActiveExtensionsResponse, error) {
	data := new(GetUserActiveExtensionsResponse)
	_, err := client.sendRequest(getUsersActiveExtensionsPath, opt, data, "GET")
	if err != nil {
		return nil, err
	}
	return data, err
}
