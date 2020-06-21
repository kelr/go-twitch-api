package pubsub

const (
	whispersTopic = "whispers."
)

// WhispersEvent contains the type and data payload for a mod action event
type WhispersEvent struct {
	Type string       `json:"type"`
	Data WhispersData `json:"data"`
}

type whispersEventDecode struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

// WhispersData contains information about a whisper
type WhispersData struct {
	ID        int               `json:"id"`
	MessageID string            `json:"message_id"`
	ThreadID  string            `json:"thread_id"`
	Body      string            `json:"body"`
	SentTs    int               `json:"sent_ts"`
	FromID    int               `json:"from_id"`
	Tags      WhispersTags      `json:"tags"`
	Recipient WhispersRecipient `json:"recipient"`
	Nonce     string            `json:"nonce"`
}

// WhispersTags represents information about the sender of the whisper
type WhispersTags struct {
	Login       string           `json:"login"`
	DisplayName string           `json:"display_name"`
	Color       string           `json:"color"`
	Emotes      []string         `json:"emotes"`
	Badges      []WhispersBadges `json:"badges"`
}

// WhispersBadges represents the badges shown by the sender of the whisper
type WhispersBadges struct {
	ID      string `json:"id"`
	Version string `json:"version"`
}

// WhispersRecipient represents the recipient of the whisper
type WhispersRecipient struct {
	ID           int     `json:"id"`
	Username     string  `json:"username"`
	DisplayName  string  `json:"display_name"`
	Color        string  `json:"color"`
	ProfileImage *string `json:"profile_image"`
}

// ListenWhispers subscribes a handler function to the Whispers topic with the provided id.
// The handler will be called with a populated WhispersData struct when the event is received.
func (c *Client) ListenWhispers(handler func(*WhispersData)) {
	c.whispersHandler = handler
	if c.IsConnected() {
		c.listen(&[]string{whispersTopic + c.ID})
	}
}

// UnlistenWhispers removes the current handler function from the Whispers event topic and
// unlistens from the topic.
func (c *Client) UnlistenWhispers() {
	c.whispersHandler = nil
	if c.IsConnected() {
		c.unlisten(&[]string{whispersTopic + c.ID})
	}
}
