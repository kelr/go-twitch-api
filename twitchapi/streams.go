package twitchapi

const (
	getStreamsPath         = "/streams"
	getStreamsMetadataPath = "/streams/metadata"
)

// Defines the options available for Get Streams
type GetStreamsOpt struct {
	After     string `json:"after,omitempty"`
	Before    string `json:"before,omitempty"`
	First     int    `json:"first,omitempty"`
	GameID    string `json:"game_id,omitempty"`
	Language  string `json:"language,omitempty"`
	UserID    string `json:"user_id,omitempty"`
	UserLogin string `json:"user_login,omitempty"`
}

// Response structure for a Get Streams command
type GetStreamsResponse struct {
	Data []struct {
		GameID       string `json:"game_id,omitempty"`
		ID           string `json:"id,omitempty"`
		Language     string `json:"language,omitempty"`
		StartedAt    string `json:"started_at,omitempty"`
		TagIDs       string `json:"tag_ids,omitempty"`
		ThumbnailURL string `json:"thumbnail_url,omitempty"`
		Title        string `json:"title,omitempty"`
		Type         string `json:"type,omitempty"`
		UserID       string `json:"user_id,omitempty"`
		Username     string `json:"user_name,omitempty"`
		ViewerCount  int    `json:"viewer_count,omitempty"`
	} `json:"data,omitempty"`
	Pagination struct {
		Cursor string `json:"cursor,omitempty"`
	} `json:"pagination,omitempty"`
}

// Return a slice representing the top active streams sorted by viewcount. Also
// returns a Pagination field used to query for more streams
//
// https://dev.twitch.tv/docs/api/reference#get-streams
func (client *TwitchClient) GetStreams(opt *GetStreamsOpt) (*GetStreamsResponse, error) {
	if opt != nil {
	}
	data := new(GetStreamsResponse)
	_, err := client.sendRequest(getStreamsPath, nil, data)
	if err != nil {
		return nil, err
	}
	return data, err
}
