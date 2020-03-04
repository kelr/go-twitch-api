package twitchapi

const (
	getUsersPath = "/users"
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
	_, err := client.sendRequest(getUsersPath, opt, data)
	if err != nil {
		return nil, err
	}
	return data, err
}
