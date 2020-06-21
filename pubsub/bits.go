package pubsub

import (
	"time"
)

const (
	bitsTopic      = "channel-bits-events-v2."
	bitsBadgeTopic = "channel-bits-badge-unlocks."
)

// BitsEvent contains the type and data payload for a bits event.
type BitsEvent struct {
	Data        BitsData `json:"data"`
	Version     string   `json:"version"`
	MessageType string   `json:"message_type"`
	MessageID   string   `json:"message_id"`
}

// BitsData contains information about a bits event.
type BitsData struct {
	Username         string               `json:"user_name"`
	ChannelName      string               `json:"channel_name"`
	UserID           string               `json:"user_id"`
	ChannelID        string               `json:"channel_id"`
	Time             time.Time            `json:"time"`
	ChatMessage      string               `json:"chat_message"`
	BitsUsed         int                  `json:"bits_used"`
	TotalBitsUsed    int                  `json:"total_bits_used"`
	IsAnonymous      bool                 `json:"is_anonymous"`
	Context          string               `json:"context"`
	BadgeEntitlement BadgeEntitlementData `json:"badge_entitlement"`
}

// BadgeEntitlementData represents entitlement information inside a Bits event.
type BadgeEntitlementData struct {
	NewVersion      int `json:"new_version"`
	PreviousVersion int `json:"previous_version"`
}

// ListenBits subscribes a handler function to the Bits event topic with the provided id.
// The handler will be called with a populated BitsData struct when the event is received.
func (c *Client) ListenBits(handler func(*BitsData)) {
	c.bitsHandler = handler
	if c.IsConnected() {
		c.listen(&[]string{bitsTopic + c.ID})
	}
}

// UnlistenBits removes the current handler function from the Bits event topic and
// unlistens from the topic.
func (c *Client) UnlistenBits() {
	c.bitsHandler = nil
	if c.IsConnected() {
		c.unlisten(&[]string{bitsTopic + c.ID})
	}
}

// BitsBadgeEvent contains the type and data payload for a bits badge event.
type BitsBadgeEvent struct {
	Type string        `json:"type"`
	Data BitsBadgeData `json:"data"`
}

// BitsBadgeData contains information about a bits badge event.
type BitsBadgeData struct {
}

// ListenBitsBadge subscribes a handler function to the Bits badge event topic with the provided id.
// The handler will be called with a populated BitsBadgeData struct when the event is received.
func (c *Client) ListenBitsBadge(handler func(*BitsBadgeData)) {
	c.bitsBadgeHandler = handler
	if c.IsConnected() {
		c.listen(&[]string{bitsBadgeTopic + c.ID})
	}
}

// UnlistenBitsBadge removes the current handler function from the Bits badge event topic and
// unlistens from the topic.
func (c *Client) UnlistenBitsBadge() {
	c.bitsBadgeHandler = nil
	if c.IsConnected() {
		c.unlisten(&[]string{bitsBadgeTopic + c.ID})
	}
}
