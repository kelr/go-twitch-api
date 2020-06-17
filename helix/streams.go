package helix

import (
	"encoding/json"
)

const (
	getStreamsPath         = "/streams"
	getStreamsMetadataPath = "/streams/metadata"
)

// GetStreamsOpt defines the options available for Get Streams.
type GetStreamsOpt struct {
	After     string `url:"after,omitempty"`
	Before    string `url:"before,omitempty"`
	First     int    `url:"first,omitempty"`
	GameID    string `url:"game_id,omitempty"`
	Language  string `url:"language,omitempty"`
	UserID    string `url:"user_id,omitempty"`
	UserLogin string `url:"user_login,omitempty"`
}

// GetStreamsResponse represents a response from a Get Streams command.
type GetStreamsResponse struct {
	Data []struct {
		GameID       string `json:"game_id,omitempty"`
		ID           string `json:"id,omitempty"`
		Language     string `json:"language,omitempty"`
		StartedAt    string `json:"started_at,omitempty"`
		ThumbnailURL string `json:"thumbnail_url,omitempty"`
		Title        string `json:"title,omitempty"`
		Type         string `json:"type,omitempty"`
		UserID       string `json:"user_id,omitempty"`
		Username     string `json:"user_name,omitempty"`
		ViewerCount  int    `json:"viewer_count,omitempty"`
		TagIDs       []string `json:"tag_ids,omitempty"`
	} `json:"data,omitempty"`

	Pagination struct {
		Cursor string `json:"cursor,omitempty"`
	} `json:"pagination,omitempty"`
}

// GetStreams returns a slice representing the top active streams sorted by viewcount. Also
// returns a Pagination field used to query for more streams
//
// https://dev.twitch.tv/docs/api/reference#get-streams
func (client *Client) GetStreams(opt *GetStreamsOpt) (*GetStreamsResponse, error) {
	data := new(GetStreamsResponse)
	resp, err := client.getRequest(getStreamsPath, opt)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resp.Data, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
