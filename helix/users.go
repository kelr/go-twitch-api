package helix

import (
	"time"
	"encoding/json"
)

const (
	getUsersPath                 = "/users"
	getUsersFollowsPath          = "/users/follows"
	getUsersExtensionsPath       = "/users/extensions/list"
	getUsersActiveExtensionsPath = "/users/extensions"
)

// PaginationData represents the current ID for a multi-page response.
type PaginationData struct {
	Cursor string `json:"cursor,omitempty"`
}

// GetUsersOpt defines the options available for Get Users.
type GetUsersOpt struct {
	ID    string `url:"id,omitempty"`
	Login string `url:"login,omitempty"`
}

// GetUsersData represents information about a user.
type GetUsersData struct {
	ID              string `json:"id,omitempty"`
	Login           string `json:"login,omitempty"`
	DisplayName     string `json:"display_name,omitempty"`
	Type            string `json:"type,omitempty"`
	BroadcasterType string `json:"broadcaster_type,omitempty"`
	Description     string `json:"description,omitempty"`
	ProfileImageURL string `json:"profile_image_url,omitempty"`
	OfflineImageURL string `json:"offline_image_url,omitempty"`
	ViewCount       int `json:"view_count,omitempty"`
	Email           string `json:"email,omitempty"`
}

// GetUsersResponse represents a response from a Get Users command
type GetUsersResponse struct {
	Data []GetUsersData `json:"data,omitempty"`
}

// GetUsers returns information about a Twitch user.
// Returns a GetUsersResponse constructed from the response from the API endpoint.
//
// https://dev.twitch.tv/docs/api/reference#get-users
func (client *Client) GetUsers(opt *GetUsersOpt) (*GetUsersResponse, error) {
	data := new(GetUsersResponse)

	resp, err := client.getRequest(getUsersPath, opt)
	if err != nil {
		return nil, err
	}

	// Decode the response
	err = json.Unmarshal(resp.Data, data)
	if err != nil {
		return nil, err
	}
	
	return data, nil
}

// GetUsersFollowsOpt defines the options available for Get Users Follows.
type GetUsersFollowsOpt struct {
	After  string `url:"after,omitempty"`
	First  int    `url:"first,omitempty"`
	FromID string `url:"from_id,omitempty"`
	ToID   string `url:"to_id,omitempty"`
}

// GetUsersFollowsData represents information about a user follow.
type GetUsersFollowsData struct {
	FollowedAt time.Time `json:"followed_at,omitempty"`
	FromID     string    `json:"from_id,omitempty"`
	FromName   string    `json:"from_name,omitempty"`
	ToID       string    `json:"to_id,omitempty"`
	ToName     string    `json:"to_name,omitempty"`
}

// GetUsersFollowsResponse represents a response from a Get Users Follows command
type GetUsersFollowsResponse struct {
	Total      int                   `json:"total,omitempty"`
	Data       []GetUsersFollowsData `json:"data,omitempty"`
	Pagination PaginationData        `json:"pagination,omitempty"`
}

// GetUsersFollows obtains information about who a user is following or who follows a user.
// Returns a GetUsersFollowsResponse constructed from the response from the API endpoint.
//
// https://dev.twitch.tv/docs/api/reference#get-users-follows
func (client *Client) GetUsersFollows(opt *GetUsersFollowsOpt) (*GetUsersFollowsResponse, error) {
	data := new(GetUsersFollowsResponse)
	resp, err := client.getRequest(getUsersFollowsPath, opt)
	if err != nil {
		return nil, err
	}

	// Decode the response
	err = json.Unmarshal(resp.Data, data)
	if err != nil {
		return nil, err
	}
	
	return data, nil
}

// UpdateUserOpt defines the options available for Update User
type UpdateUserOpt struct {
	Description string `url:"description"`
}

// UpdateUser updates the description of a user. Requires a user token for the user to be updated.
// Returns a GetUsersResponse constructed from the response from the API endpoint.
// Requires scope: user:edit
//
// https://dev.twitch.tv/docs/api/reference#update-user
func (client *Client) UpdateUser(opt *UpdateUserOpt) (*GetUsersResponse, error) {
	data := new(GetUsersResponse)

	resp, err := client.putRequest(getUsersPath, opt)
	if err != nil {
		return nil, err
	}

	// Decode the response
	err = json.Unmarshal(resp.Data, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetUsersExtensionsData represents information about a users extensions.
type GetUsersExtensionsData struct {
	ID          string   `json:"id,omitempty"`
	Version     string   `json:"version,omitempty"`
	Name        string   `json:"name,omitempty"`
	CanActivate bool     `json:"can_activate,omitempty"`
	Type        []string `json:"type,omitempty"`
}

// GetUserExtensionsResponse represents a response from a Get Users Extensions command
type GetUserExtensionsResponse struct {
	Data []GetUsersExtensionsData `json:"data,omitempty"`
}

// GetUserExtensions returns information about active and inactive extensions for a user identified by the user token.
// Requires scope user:read:broadcast
//
// https://dev.twitch.tv/docs/api/reference#get-users-follows
func (client *Client) GetUserExtensions() (*GetUserExtensionsResponse, error) {
	data := new(GetUserExtensionsResponse)

	resp, err := client.getRequest(getUsersExtensionsPath, nil)
	if err != nil {
		return nil, err
	}

	// Decode the response
	err = json.Unmarshal(resp.Data, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetUserActiveExtensionsOpt defines the options available for Get User Active Extensions
type GetUserActiveExtensionsOpt struct {
	UserID string `url:"user_id"`
}

// ActiveExtension represents a currently active extension.
type ActiveExtension struct {
	Active  bool   `json:"active"`
	ID      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
	X       int    `json:"x,omitempty"`
	Y       int    `json:"y,omitempty"`
}

// GetUserActiveExtensionsData represents information about a users active extensions.
type GetUserActiveExtensionsData struct {
	Component map[string]ActiveExtension `json:"component,omitempty"`
	Overlay   map[string]ActiveExtension `json:"overlay,omitempty"`
	Panel     map[string]ActiveExtension `json:"panel,omitempty"`
}

// GetUserActiveExtensionsResponse represents a response from a Get Users Extensions command
type GetUserActiveExtensionsResponse struct {
	Data GetUserActiveExtensionsData `json:"data,omitempty"`
}

// GetUserActiveExtensions returns information about active and inactive extensions for a user identified by the user token
// Requires scope user:read:broadcast or user:edit:broadcast
//
// https://dev.twitch.tv/docs/api/reference#get-user-active-extensions
func (client *Client) GetUserActiveExtensions(opt *GetUserActiveExtensionsOpt) (*GetUserActiveExtensionsResponse, error) {
	data := new(GetUserActiveExtensionsResponse)
	resp, err := client.getRequest(getUsersActiveExtensionsPath, opt)
	if err != nil {
		return nil, err
	}
	// Decode the response
	err = json.Unmarshal(resp.Data, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
