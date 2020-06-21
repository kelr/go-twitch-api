package pubsub

import (
	"encoding/json"
	"fmt"
	"strings"
)

type pubSubEvent struct {
	Type string     `json:"type,omitempty"`
	Data pubSubData `json:"data,omitempty"`
}

type pubSubData struct {
	Topic   string `json:"topic,omitempty"`
	Message string `json:"message,omitempty"`
}

type pubSubResponse struct {
	Type  string `json:"type"`
	Nonce string `json:"nonce"`
	Error string `json:"error"`
}

// handle determines the type of message and calls the corresponding handler.
func (c *Client) handle(msg []byte) {
	builtMsg := new(pubSubEvent)
	err := json.Unmarshal(msg, builtMsg)
	if err != nil {
		fmt.Println(msg, err)
		return
	}
	// Since the Message field is a string of encoded JSON, we have to
	// determine the type of message and unmarshal the Message field specifically.
	switch builtMsg.Type {
	case "MESSAGE":
		s := strings.SplitAfter(builtMsg.Data.Topic, ".")
		topic := s[0]
		switch topic {
		case channelPointTopic:
			if c.channelPointHandler != nil {
				if err = c.handleChannelPointsEvent(builtMsg.Data.Message); err != nil {
					fmt.Println(err)
				}
			}
		case chatModActionsTopic:
			if c.chatModActionsHandler != nil {
				if err = c.handleChatModActionsEvent(builtMsg.Data.Message); err != nil {
					fmt.Println(err)
				}
			}
		case whispersTopic:
			if c.whispersHandler != nil {
				if err = c.handleWhispersEvent(builtMsg.Data.Message); err != nil {
					fmt.Println(err)
				}
			}
		case subsTopic:
			if c.subsHandler != nil {
				if err = c.handleSubsEvent(builtMsg.Data.Message); err != nil {
					fmt.Println(err)
				}
			}
		case bitsTopic:
			if c.bitsHandler != nil {
				if err = c.handleBitsEvent(builtMsg.Data.Message); err != nil {
					fmt.Println(err)
				}
			}
		case bitsBadgeTopic:
			if c.bitsBadgeHandler != nil {
				if err = c.handleBitsBadgeEvent(builtMsg.Data.Message); err != nil {
					fmt.Println(err)
				}
			}
		default:
			fmt.Println("Unknown topic:", topic)
		}
	case "RESPONSE":
		if err = c.handleResponse(msg); err != nil {
			fmt.Println(err)
		}
	case "PONG":
		c.handlePong()
	case "RECONNECT":
		c.handleReconnect()
	default:
		fmt.Println("PubSub Client unknown message received:", builtMsg.Type)
	}
}

func (c *Client) handleChannelPointsEvent(message string) error {
	event := new(ChannelPointsEvent)
	err := json.Unmarshal([]byte(message), event)
	if err != nil {
		return err
	}
	c.channelPointHandler(&event.Data)
	return nil
}

func (c *Client) handleChatModActionsEvent(message string) error {
	event := new(ChatModActionsEvent)
	err := json.Unmarshal([]byte(message), event)
	if err != nil {
		return err
	}
	c.chatModActionsHandler(&event.Data)
	return nil
}

func (c *Client) handleWhispersEvent(message string) error {
	event := new(WhispersEvent)
	tmp := new(whispersEventDecode)
	err := json.Unmarshal([]byte(message), tmp)
	if err != nil {
		return err
	}

	switch tmp.Type {
	case "whisper_sent":
		// TODO
	case "whisper_received":
		// TODO
	case "thread":
		// Thread is not handled as of yet.
		return nil
	default:
		return fmt.Errorf("Unknown whispers type: %s", tmp.Type)
	}

	// The data field is a double escaped JSON string so we need to unmarshal it twice.
	event.Type = tmp.Type
	err = json.Unmarshal([]byte(tmp.Data), &event.Data)
	if err != nil {
		return err
	}
	c.whispersHandler(&event.Data)
	return nil
}

func (c *Client) handleSubsEvent(message string) error {
	event := new(SubsData)
	err := json.Unmarshal([]byte(message), event)
	if err != nil {
		return err
	}
	c.subsHandler(event)
	return nil
}

func (c *Client) handleBitsEvent(message string) error {
	event := new(BitsEvent)
	err := json.Unmarshal([]byte(message), event)
	if err != nil {
		return err
	}
	c.bitsHandler(&event.Data)
	return nil
}

func (c *Client) handleBitsBadgeEvent(message string) error {
	event := new(BitsBadgeEvent)
	err := json.Unmarshal([]byte(message), event)
	if err != nil {
		return err
	}
	c.bitsBadgeHandler(&event.Data)
	return nil
}

// handleResponse checks for errors in the RESPONSE message received after a LISTEN request.
func (c *Client) handleResponse(message []byte) error {
	resp := new(pubSubResponse)
	err := json.Unmarshal(message, resp)
	if err != nil {
		fmt.Println(message, err)
		return err
	}
	if resp.Error != "" {
		return fmt.Errorf("PubSub client received error response: %s", resp.Error)
	}
	return nil
}

// handlePong notifies on the pongRx channel that a Pong was received.
func (c *Client) handlePong() {
	c.pongRx <- true
}

// handlReconnect prepares for the PubSub endpoint to go down within the next 30s.
func (c *Client) handleReconnect() {
	// TODO: prepare for reconnect after shutdown within 30s
	fmt.Println("PubSub client received reconnect message.")
}
