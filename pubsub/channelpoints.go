package pubsub

import (
	"time"
	"fmt"
)

const (
	channelPointTopic = "channel-points-channel-v1."
)

// ChannelPointsEvent contains the type and data payload for a channel points event
type ChannelPointsEvent struct {
	Type string            `json:"type"`
	Data ChannelPointsData `json:"data"`
}

// ChannelPointsData contains the time the reward was redeemed and the redemption data
type ChannelPointsData struct {
	TimeStamp  time.Time      `json:"timestamp"`
	Redemption RedemptionData `json:"redemption"`
}

// RedemptionData contains metadata about the redeemed reward
type RedemptionData struct {
	Id         string           `json:"id"`
	User       RedemptionUser   `json:"user"`
	ChannelId  string           `json:"channel_id"`
	RedeemedAt time.Time        `json:"redeemed_at"`
	Reward     RedemptionReward `json:"reward"`
	UserInput  string           `json:"user_input"`
	Status     string           `json:"status"`
}

// RedemptionUser represents the user who redeemed the reward
type RedemptionUser struct {
	Id          string `json:"id"`
	Login       string `json:"login"`
	DisplayName string `json:"display_name"`
}

// RedemptionReward represents information about the reward redeemed
type RedemptionReward struct {
	Id                    string                 `json:"id"`
	ChannelId             string                 `json:"channel_id"`
	Title                 string                 `json:"title"`
	Prompt                string                 `json:"prompt"`
	Cost                  int                    `json:"cost"`
	IsUserInputRequired   bool                   `json:"is_user_input_required"`
	IsSubOnly             bool                   `json:"is_sub_only"`
	Image                 RedemptionImage        `json:"image"`
	DefaultImage          RedemptionImage        `json:"default_image"`
	BackgroundColor       string                 `json:"background_color"`
	IsEnabled             bool                   `json:"is_enabled"`
	IsPaused              bool                   `json:"is_paused"`
	IsInStock             bool                   `json:"is_in_stock"`
	MaxPerStream          RedemptionMaxPerStream `json:"max_per_stream"`
	ShouldRedemptionsSkip bool                   `json:"should_redemptions_skip_request_queue"`
	TemplateId            string                 `json:"template_id"`
	UpdatedForIndicatorAt time.Time              `json:"updated_for_indicator_at"`
}

// RedemptionImage represents the cute image used on the redemption button
type RedemptionImage struct {
	Url1x string `json:"url_1x"`
	Url2x string `json:"url_2x"`
	Url4x string `json:"url_4x"`
}

// RedemptionMaxPerStream represents information about redemption limits per stream
type RedemptionMaxPerStream struct {
	IsEnabled    bool `json:"is_enabled"`
	MaxPerStream int  `json:"max_per_stream"`
}

func (c *PubSubClient) ListenChannelPoints(id string, handler func(*ChannelPointsEvent)) error {
	if _, ok := c.channelPointHandlers[id]; !ok {
		c.channelPointHandlers[id] = handler
		if c.IsConnected() {
			c.listen(&[]string{channelPointTopic + id})
		}
	} else {
		return fmt.Errorf("Channel Points handler already registered for Id: %s", id)
	}
	return nil
}