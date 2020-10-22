package helix

import (
	"encoding/json"
)

const (
	getClipsPath = "/clips"
)

// GetClipsOpt defines the options available for Get Clips.
type GetClipsOpt struct {
	ID            string `url:"id,omitempty"`
	BroadcasterID string `url:"broadcaster_id,omitempty"`
	GameID        string `url:"game_id,omitempty"`
	After         string `url:"after,omitempty"`
	First         int    `url:"first,omitempty"`
	StartedAt     string `url:"started_at,omitempty"`
	EndedAt       string `url:"ended_at,omitempty"`
}

// GetClipsData represents metadata about a clip.
type GetClipsData struct {
	ID              string `json:"id,omitempty"`
	URL             string `json:"url,omitempty"`
	EmbedURL        string `json:"embed_url,omitempty"`
	BroadcasterID   string `json:"broadcaster_id,omitempty"`
	BroadcasterName string `json:"broadcaster_name,omitempty"`
	CreatorID       string `json:"creator_id,omitempty"`
	CreatorName     string `json:"creator_name,omitempty"`
	VideoID         string `json:"video_id,omitempty"`
	GameID          string `json:"game_id,omitempty"`
	Language        string `json:"language,omitempty"`
	Title           string `json:"title,omitempty"`
	ViewCount       int    `json:"view_count,omitempty"`
	CreatedAt       string `json:"created_at,omitempty"`
	ThumbnailURL    string `json:"thumbnail_url,omitempty"`
}

// GetClipsResponse represents the response from a Get Clips command.
type GetClipsResponse struct {
	Data       []GetClipsData `json:"data,omitempty"`
	Pagination PaginationData
}

// GetClips gets information by clip id, broadcaster id or game id.
//
// https://dev.twitch.tv/docs/api/reference/#get-clips
func (client *Client) GetClips(opt *GetClipsOpt) (*GetClipsResponse, error) {
	data := new(GetClipsResponse)
	resp, err := client.getRequest(getClipsPath, opt)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(resp.Data, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
