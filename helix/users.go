package helix //import "github.com/kelr/go-twitch-api/helix"

const (
	getUsersPath        = "/users"
	getUsersFollowsPath = "/users/follows"
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
		FollowedAt string `json:"followed_at,omitempty"`
		FromID     string `json:"from_id,omitempty"`
		FromName   string `json:"from_name,omitempty"`
		ToID       string `json:"to_id,omitempty"`
		ToName     string `json:"to_name,omitempty"`
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

// Defines the options available for Get Users Follows
type UpdateUserOpt struct {
	Description  string `url:"description"`
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