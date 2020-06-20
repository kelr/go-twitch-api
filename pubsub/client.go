// Package pubsub provides a client used to subscribe to updates from the Twitch PubSub endpoints.
package pubsub

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"golang.org/x/oauth2"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	pubSubURL        = "wss://pubsub-edge.twitch.tv"
	nonceSet         = "abcdedfghijklmnopqrstABCDEFGHIJKLMNOP123456789"
	pingPeriod       = 250 * time.Second
	pongWaitPeriod   = 10 * time.Second
	pingMsg          = `{"type": "PING"}`
	maxReconnectTime = 600
)

// Message contains the entire message structure receieved from a topic
type Message struct {
	Type string     `json:"type,omitempty"`
	Data Data `json:"data,omitempty"`
}

// Data contains the topic and the message payload as an encoded JSON string.
type Data struct {
	Topic   string `json:"topic,omitempty"`
	Message string `json:"message,omitempty"`
}

type pubSubResponse struct {
	Type  string `json:"type"`
	Nonce string `json:"nonce"`
	Error string `json:"error"`
}

type pubSubListenRequest struct {
	Type  string `json:"type"`
	Nonce string `json:"nonce"`
	Data  struct {
		Topics    []string `json:"topics"`
		AuthToken string   `json:"auth_token"`
	} `json:"data"`
}

// Client represents a connection and its state to the Twitch pubsub endpoint.
type Client struct {
	conn                   *websocket.Conn
	refreshClient          *http.Client
	sendChan               chan []byte
	AuthToken              *oauth2.Token
	stop                   chan bool
	pongRx                 chan bool
	reconnectTime          int
	isConnected            bool
	mu                     *sync.Mutex
	channelPointHandlers   map[string]func(*ChannelPointsEvent)
	chatModActionsHandlers map[string]func(*ChatModActionsEvent)
	whispersHandlers       map[string]func(*WhispersEvent)
}

// NewClient returns a new Client to communicate with the PubSub endpoints.
func NewClient(config *oauth2.Config, userToken *oauth2.Token) *Client {
	return &Client{
		conn:                   nil,
		AuthToken:              userToken,
		refreshClient:          config.Client(context.Background(), userToken),
		sendChan:               make(chan []byte, 256),
		stop:                   make(chan bool),
		pongRx:                 make(chan bool, 1),
		reconnectTime:          1,
		isConnected:            false,
		mu:                     &sync.Mutex{},
		channelPointHandlers:   make(map[string]func(*ChannelPointsEvent)),
		chatModActionsHandlers: make(map[string]func(*ChatModActionsEvent)),
		whispersHandlers:       make(map[string]func(*WhispersEvent)),
	}
}

// IsConnected is a thread-safe check of whether or not the Client is connected.
func (c *Client) IsConnected() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	result := c.isConnected
	return result
}

// Connect to the Twitch PubSub endpoint and listen on all registered topics.
// Will automatically reconnect on failure.
// with exponential backoff. Returns an error if the client is already connected.
func (c *Client) Connect() error {
	if !c.IsConnected() {
		conn, _, err := websocket.DefaultDialer.Dial(pubSubURL, nil)
		if err != nil {
			go c.reconnect()
			return err
		}
		c.conn = conn
		go c.reader()
		go c.writer()

		c.listenAll()

		c.mu.Lock()
		c.isConnected = true
		c.mu.Unlock()
		fmt.Println("PubSub Client connected")
	} else {
		errors.New("PubSub Client is already connected")
	}
	return nil
}

func (c *Client) listenAll() {
	// Listen on all registered topics
	var topics []string
	for id := range c.channelPointHandlers {
		fmt.Println("Listening:", channelPointTopic+id)
		topics = append(topics, channelPointTopic+id)
	}
	for id := range c.chatModActionsHandlers {
		fmt.Println("Listening:", chatModActionsTopic+id)
		topics = append(topics, chatModActionsTopic+id)
	}
	for id := range c.whispersHandlers {
		fmt.Println("Listening:", whispersTopic+id)
		topics = append(topics, whispersTopic+id)
	}
	if len(topics) > 0 {
		c.listen(&topics)
	}
}

// Close disconnects the client from the Twitch PubSub endpoint.
// If the client is already connected, Close() will return an error.
func (c *Client) Close() error {
	if c.IsConnected() {
		c.mu.Lock()
		c.isConnected = false
		c.mu.Unlock()
		close(c.stop)
		c.conn.Close()
	} else {
		return errors.New("PubSub Client connection is already closed")
	}
	return nil
}

// Shutdown tells every goroutine to stop if any one of them signals shutdown.
// Closes the connection and attempt to reconnect.
func (c *Client) shutdown() {
	c.Close()
	go c.reconnect()
}

// Updates the exponential backoff reconnect time and attempts to reconnect
// to the Twitch PubSub endpoint after this time period.
func (c *Client) reconnect() {
	if c.reconnectTime < maxReconnectTime {
		c.reconnectTime *= 2
	}
	fmt.Println("PubSub Client lost connection, retrying in:", c.reconnectTime, "seconds")
	reconnectTimer := time.NewTimer(time.Duration(c.reconnectTime) * time.Second)
	<-reconnectTimer.C
	c.stop = make(chan bool)
	c.Connect()
}

// reader reads messages from the connection and passes it onto the handler.
func (c *Client) reader() {
	for {
		select {
		case <-c.stop:
			return
		default:
			_, msg, err := c.conn.ReadMessage()
			if err != nil {
				fmt.Println("PubSub error in rx:", err)
				c.shutdown()
				return
			}
			c.handle(msg)
		}
	}
}

// handle determines the type of message and calls the corresponding handler.
func (c *Client) handle(msg []byte) {
	builtMsg := new(Message)
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
		topic, id := s[0], s[1]
		switch topic {
		case channelPointTopic:
			if err = c.handleChannelPointsEvent(builtMsg.Data.Message, id); err != nil {
				fmt.Println(err)
				return
			}
		case chatModActionsTopic:
			if err = c.handleChatModActionsEvent(builtMsg.Data.Message, id); err != nil {
				fmt.Println(err)
				return
			}
		case whispersTopic:
			if err = c.handleWhispersEvent(builtMsg.Data.Message, id); err != nil {
				fmt.Println(err)
				return
			}
		default:
			fmt.Println("Unknown topic:", topic)
		}
	case "RESPONSE":
		if err = c.handleResponse(msg); err != nil {
			fmt.Println(err)
			return
		}
	case "PONG":
		c.handlePong()
	case "RECONNECT":
		c.handleReconnect()
	default:
		fmt.Println("PubSub Client unknown message received:", builtMsg.Type)
	}
}

func (c *Client) handleChannelPointsEvent(message string, id string) error {
	event := new(ChannelPointsEvent)
	err := json.Unmarshal([]byte(message), event)
	if err != nil {
		return err
	}
	c.channelPointHandlers[id](event)
	return nil
}

func (c *Client) handleChatModActionsEvent(message string, id string) error {
	event := new(ChatModActionsEvent)
	err := json.Unmarshal([]byte(message), event)
	if err != nil {
		return err
	}
	c.chatModActionsHandlers[id](event)
	return nil
}

func (c *Client) handleWhispersEvent(message string, id string) error {
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
	c.whispersHandlers[id](event)
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

// writer handles transmitting regular ping messages and determines if a pong response is in time.
// Also writes any messages from the send channel.
func (c *Client) writer() {
	pingTicker := time.NewTicker(pingPeriod)
	defer pingTicker.Stop()
	for {
		select {
		case <-c.stop:
			return
		case msg := <-c.sendChan:
			if err := c.write(msg); err != nil {
				c.shutdown()
				return
			}
		case <-pingTicker.C:
			if err := c.write([]byte(pingMsg)); err != nil {
				c.shutdown()
				return
			}
			go func() {
				pongWait := time.NewTimer(pongWaitPeriod)
				select {
				case <-pongWait.C:
					pongWait.Stop()
					fmt.Println("PubSub server did not respond to ping in time, reconnecting")
					c.shutdown()
					return
				case <-c.pongRx:
					return
				}
			}()
		}
	}
}

// write calls WriteMessage on the underlying client.
func (c *Client) write(msg []byte) error {
	err := c.conn.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		fmt.Println("PubSub error in tx:", err)
		return err
	}
	return nil
}

// listen creates a listen request and sends it to the send channel.
func (c *Client) listen(topics *[]string) {
	request := pubSubListenRequest{
		Type:  "LISTEN",
		Nonce: generateNonce(15),
	}
	request.Data.Topics = *topics
	request.Data.AuthToken = c.AuthToken.AccessToken
	bytes, _ := json.Marshal(request)
	c.sendChan <- bytes
}

// generateNonce creates a nonce string of variable length.
func generateNonce(length int) string {
	var curr strings.Builder
	for i := 0; i < length; i++ {
		curr.WriteString(string(nonceSet[rand.Intn(len(nonceSet))]))
	}
	return curr.String()
}
