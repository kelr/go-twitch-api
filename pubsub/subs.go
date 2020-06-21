package pubsub

const (
	subsTopic = "channel-subscribe-events-v1."
)

// SubsEvent contains the type and data payload for a subscription event.
type SubsEvent struct {
	Type string   `json:"type"`
	Data SubsData `json:"data"`
}

// SubsData contains information about a subscription event.
type SubsData struct {
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
