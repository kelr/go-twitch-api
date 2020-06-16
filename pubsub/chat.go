package pubsub

import (
	"fmt"
)

const (
	chatModActionsTopic = "chat_moderator_actions."
)

// ChatModActionsEvent contains the type and data payload for a mod action event
type ChatModActionsEvent struct {
	Data struct {
		Type string `json:"type"`
		ModerationAction string `json:"moderation_action"`
		Args []string `json:"args"`
		CreatedBy string `json:"created_by"`
		CreatedByUserId string `json:"created_by_user_id"`
		MsgId string `json:"msg_id"`
		TargetUserId string `json:"target_user_id"`
		TargetUserLogin string `json:"target_user_login"`
		FromAutomod bool `json:"from_automod"`
	} `json:"data"`
}


func (c *PubSubClient) ListenChatModActions(id string, handler func(*ChatModActionsEvent)) error {
	if _, ok := c.chatModActionsHandlers[id]; !ok {
		c.chatModActionsHandlers[id] = handler
		if c.IsConnected() {
			c.listen(&[]string{chatModActionsTopic + id})
		}
	} else {
		return fmt.Errorf("Chat Mod Actions handler already registered for Id: %s", id)
	}
	return nil
}