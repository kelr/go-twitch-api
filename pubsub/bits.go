package pubsub

const (
	bitsTopic      = "channel-bits-events-v2."
	bitsBadgeTopic = "channel-bits-badge-unlocks."
)

// BitsEvent contains the type and data payload for a bits event.
type BitsEvent struct {
	Type string   `json:"type"`
	Data BitsData `json:"data"`
}

// BitsData contains information about a bits event.
type BitsData struct {
}

// ListenBits subscribes a handler function to the Bits event topic with the provided id.
// The handler will be called with a populated BitsData struct when the event is received.
func (c *Client) ListenBits(handler func(*BitsData)) {
	c.bitsHandler = handler
	if c.IsConnected() {
		c.listen(&[]string{bitsTopic + c.ID})
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
