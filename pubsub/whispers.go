package pubsub

import (
	"fmt"
)

const (
	whispersTopic = "whispers."
)

// ChatModActionsEvent contains the type and data payload for a mod action event
type WhispersEvent struct {
	Type string       `json:"type"`
	Data WhispersData `json:"data"`
}

type WhispersEventDecode struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

type WhispersData struct {
	Id        int               `json:"id"`
	MessageId string            `json:"message_id"`
	ThreadId  string            `json:"thread_id"`
	Body      string            `json:"body"`
	SentTs    int               `json:"sent_ts"`
	FromId    int               `json:"from_id"`
	Tags      WhispersTags      `json:"tags"`
	Recipient WhispersRecipient `json:"recipient"`
	Nonce     string            `json:"nonce"`
}

type WhispersTags struct {
	Login       string           `json:"login"`
	DisplayName string           `json:"display_name"`
	Color       string           `json:"color"`
	Emotes      []string         `json:"emotes"`
	Badges      []WhispersBadges `json:"badges"`
}

type WhispersBadges struct {
	Id      string `json:"id"`
	Version string `json:"version"`
}

type WhispersRecipient struct {
	Id           int     `json:"id"`
	Username     string  `json:"username"`
	DisplayName  string  `json:"display_name"`
	Color        string  `json:"color"`
	ProfileImage *string `json:"profile_image"`
}

func (c *PubSubClient) ListenWhispers(id string, handler func(*WhispersEvent)) error {
	if _, ok := c.whispersHandlers[id]; !ok {
		c.whispersHandlers[id] = handler
		if c.IsConnected() {
			c.listen(&[]string{whispersTopic + id})
		}
	} else {
		return fmt.Errorf("Chat Mod Actions handler already registered for Id: %s", id)
	}
	return nil
}
