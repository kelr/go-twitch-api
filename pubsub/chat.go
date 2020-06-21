package pubsub

const (
	chatModActionsTopic = "chat_moderator_actions."
)

// ChatModActionsEvent contains data payload for a mod action event
type ChatModActionsEvent struct {
	Data ChatModActionsData `json:"data"`
}

// ChatModActionsData contains the data from a Chat Mod Action event
type ChatModActionsData struct {
	Type             string   `json:"type"`
	ModerationAction string   `json:"moderation_action"`
	Args             []string `json:"args"`
	CreatedBy        string   `json:"created_by"`
	CreatedByUserID  string   `json:"created_by_user_id"`
	MsgID            string   `json:"msg_id"`
	TargetUserID     string   `json:"target_user_id"`
	TargetUserLogin  string   `json:"target_user_login"`
	FromAutomod      bool     `json:"from_automod"`
}

// ListenChatModActions subscribes a handler function to the Chat Mod Actions topic with the provided id.
// The handler will be called with a populated ChatModActionsData struct when the event is received.
func (c *Client) ListenChatModActions(handler func(*ChatModActionsData)) {
	c.chatModActionsHandler = handler
	if c.IsConnected() {
		c.listen(&[]string{chatModActionsTopic + c.ID})
	}
}
