package helix

import (
	"encoding/json"
)

const (
	getVideosPath = "/videos"
)

// GetVideosOpt defines the options available for Get Videos.
type GetVideosOpt struct {
	ID       string `url:"id,omitempty"`
	UserID   string `url:"user_id,omitempty"`
	GameID   string `url:"game_id,omitempty"`
	After    string `url:"after,omitempty"`
	Before   string `url:"before,omitempty"`
	First    string `url:"first,omitempty"`
	Language string `url:"language,omitempty"`
	Period   string `url:"period,omitempty"`
	Sort     string `url:"sort,omitempty"`
	Type     string `url:"type,omitempty"`
}

// GetVideosData represents metadata about a video.
type GetVideosData struct {
	CreatedAt    string `json:"created_at,omitempty"`
	Description  string `json:"description,omitempty"`
	Duration     string `json:"duration,omitempty"`
	ID           string `json:"id,omitempty"`
	Language     string `json:"language,omitempty"`
	PublishedAt  string `json:"published_at,omitempty"`
	ThumbnailURL string `json:"thumbnail_url,omitempty"`
	Title        string `json:"title,omitempty"`
	Type         string `json:"type,omitempty"`
	URL          string `json:"url,omitempty"`
	UserID       string `json:"user_id,omitempty"`
	UserName     string `json:"user_name,omitempty"`
	ViewCount    int    `json:"view_count,omitempty"`
	Viewable     string `json:"viewable,omitempty"`
}

// GetVideosResponse represents the response from a Get Videos command.
type GetVideosResponse struct {
	Data       []GetVideosData `json:"data,omitempty"`
	Pagination PaginationData
}

// GetVideos gets information by vodep id, user id or game id.
//
// https://dev.twitch.tv/docs/api/reference/#get-videos
func (client *Client) GetVideos(opt *GetVideosOpt) (*GetVideosResponse, error) {
	data := new(GetVideosResponse)
	resp, err := client.getRequest(getVideosPath, opt)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(resp.Data, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
