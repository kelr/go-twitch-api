package helix

import (
	"encoding/json"
)

const (
	getGamesPath = "/games"
)

// GetGamesOpt defines the options available for Get Games.
type GetGamesOpt struct {
	ID   string `url:"id,omitempty"`
	Name string `url:"name,omitempty"`
}

// GetGamesData represents metadata about a game.
type GetGamesData struct {
	BoxArtURL string `json:"box_art_url,omitempty"`
	ID        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
}

// GetGamesResponse represents the response from a Get Games command.
type GetGamesResponse struct {
	Data       []GetGamesData `json:"data,omitempty"`
	Pagination PaginationData
}

// GetGames gets information by game name or game id
//
// https://dev.twitch.tv/docs/api/reference/#get-games
func (client *Client) GetGames(opt *GetGamesOpt) (*GetGamesResponse, error) {
	data := new(GetGamesResponse)
	resp, err := client.getRequest(getGamesPath, opt)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(resp.Data, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
