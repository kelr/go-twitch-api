package pubsub

import (
	"time"
)

const (
	subsTopic = "channel-subscribe-events-v1."
)

// SubsData contains the type and data payload for a subscription event.
type SubsData struct {
	BenefitEndMonth      int            `json:"benefit_end_month"`
	Username             string         `json:"user_name"`
	DisplayName          string         `json:"display_name"`
	ChannelName          string         `json:"channel_name"`
	UserID               string         `json:"user_id"`
	ChannelID            string         `json:"channel_id"`
	RecipientID          string         `json:"recipient_id"`
	RecipientUsername    string         `json:"recipient_user_name"`
	RecipientDisplayName string         `json:"recipient_display_name"`
	Time                 time.Time      `json:"time"`
	SubMessage           SubMessageData `json:"sub_message"`
	SubPlan              string         `json:"sub_plan"`
	SubPlanName          string         `json:"sub_plan_name"`
	Months               int            `json:"months"`
	Context              string         `json:"context"`
	IsGift               bool           `json:"is_gift"`
}

// SubMessageData represents data in a sub message.
type SubMessageData struct {
	Message string  `json:"message"`
	Emotes  *string `json:"emotes"`
}

// ListenSubs subscribes a handler function to the Subscriptions topic with the provided id.
// The handler will be called with a populated SubsData struct when the event is received.
func (c *Client) ListenSubs(handler func(*SubsData)) {
	c.subsHandler = handler
	if c.IsConnected() {
		c.listen(&[]string{subsTopic + c.ID})
	}
}

// UnlistenSubs removes the current handler function from the Subscriptions event topic and
// unlistens from the topic.
func (c *Client) UnlistenSubs() {
	c.subsHandler = nil
	if c.IsConnected() {
		c.unlisten(&[]string{subsTopic + c.ID})
	}
}
