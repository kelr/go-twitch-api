// Package pubsub provides a client used to subscribe to updates from the Twitch PubSub endpoints.
package pubsub

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"golang.org/x/oauth2"
	"math/rand"
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

type pubSubRequest struct {
	Type  string            `json:"type"`
	Nonce string            `json:"nonce"`
	Data  pubsubRequestData `json:"data"`
}

type pubsubRequestData struct {
	Topics    []string `json:"topics"`
	AuthToken string   `json:"auth_token"`
}

// Client represents a connection and its state to the Twitch pubsub endpoint.
type Client struct {
	AuthToken             *oauth2.Token
	ID                    string
	conn                  *websocket.Conn
	sendChan              chan []byte
	stop                  chan bool
	pongRx                chan bool
	reconnectTime         int
	isConnected           bool
	mu                    *sync.Mutex
	channelPointHandler   func(*ChannelPointsData)
	chatModActionsHandler func(*ChatModActionsData)
	whispersHandler       func(*WhispersData)
	subsHandler           func(*SubsData)
	bitsHandler           func(*BitsData)
	bitsBadgeHandler      func(*BitsBadgeData)
}

// NewClient returns a new Client to communicate with the PubSub endpoints.
func NewClient(userID string, userToken *oauth2.Token) *Client {
	return &Client{
		conn:          nil,
		ID:            userID,
		AuthToken:     userToken,
		sendChan:      make(chan []byte, 256),
		stop:          make(chan bool),
		pongRx:        make(chan bool, 1),
		reconnectTime: 1,
		isConnected:   false,
		mu:            &sync.Mutex{},
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
	if c.channelPointHandler != nil {
		topics = append(topics, channelPointTopic+c.ID)
	}
	if c.chatModActionsHandler != nil {
		topics = append(topics, chatModActionsTopic+c.ID)
	}
	if c.whispersHandler != nil {
		topics = append(topics, whispersTopic+c.ID)
	}
	if c.subsHandler != nil {
		topics = append(topics, subsTopic+c.ID)
	}
	if c.bitsHandler != nil {
		topics = append(topics, bitsTopic+c.ID)
	}
	if c.bitsBadgeHandler != nil {
		topics = append(topics, bitsBadgeTopic+c.ID)
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
			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				c.shutdown()
				return
			}
		case <-pingTicker.C:
			if err := c.conn.WriteMessage(websocket.TextMessage, []byte(pingMsg)); err != nil {
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

// listen creates a listen request and sends it to the send channel.
func (c *Client) listen(topics *[]string) {
	for _, topic := range *topics {
		fmt.Println("Listening:", topic)
	}
	request := pubSubRequest{
		Type:  "LISTEN",
		Nonce: generateNonce(15),
		Data: pubsubRequestData{
			Topics:    *topics,
			AuthToken: c.AuthToken.AccessToken,
		},
	}
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
