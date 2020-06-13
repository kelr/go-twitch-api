package helix

const (
	getClipsPath = "/clips"
)

// Defines the options available for Get Users
type GetClipsOpt struct {
	Id            string `url:"id,omitempty"`
	BroadcasterId string `url:"broadcaster_id,omitempty"`
	GameId        string `url:"game_id,omitempty"`
}

type GetClipsData struct {
	Id              string `json:"id,omitempty"`
	Url             string `json:"url,omitempty"`
	EmbedUrl        string `json:"embed_url,omitempty"`
	BroadcasterId   string `json:"broadcaster_id,omitempty"`
	BroadcasterName string `json:"braodcaster_name,omitempty"`
	CreatorId       string `json:"creator_id,omitempty"`
	CreatorName     string `json:"creator_name,omitempty"`
	VideoId         string `json:"video_id,omitempty"`
	GameId          string `json:"game_id,omitempty"`
	Language        string `json:"language,omitempty"`
	Title           string `json:"title,omitempty"`
	ViewCount       int    `json:"view_count,omitempty"`
	CreatedAt       string `json:"created_at,omitempty"`
	ThumbnailUrl    string `json:"thyumbnail_url,omitempty"`
}

// Response structure for a Get Users command
type GetClipsResponse struct {
	Data       []GetClipsData `json:"data,omitempty"`
	Pagination PaginationData
}

// Get information by clip id, broadcaster id or game id
//
// https://dev.twitch.tv/docs/api/reference/#get-clips
func (client *TwitchClient) GetClips(opt *GetClipsOpt) (*GetClipsResponse, error) {
	data := new(GetClipsResponse)
	_, err := client.sendRequest(getClipsPath, opt, data, "GET")
	if err != nil {
		return nil, err
	}
	return data, err
}
