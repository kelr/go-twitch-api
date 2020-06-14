package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"golang.org/x/oauth2"
	"math/rand"
	"net/http"
	"strings"
	"time"
	"sync"
	"errors"
)

const (
	pubSubURL        = "wss://pubsub-edge.twitch.tv"
	nonceSet         = "abcdedfghijklmnopqrstABCDEFGHIJKLMNOP123456789"
	pingPeriod       = 250 * time.Second
	pongWaitPeriod   = 10 * time.Second
	pingMsg          = `{"type": "PING"}`
	maxReconnectTime = 600
)

// PubSubMessage contains the entire message structure receieved from a topic
type PubSubMessage struct {
	Type string     `json:"type,omitempty"`
	Data PubSubData `json:"data,omitempty"`
}

// PubSubData contains the topic and the message payload as an encoded JSON string.
type PubSubData struct {
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

// PubSubClient represents a connection and its state to the Twitch pubsub endpoint.
type PubSubClient struct {
	conn          *websocket.Conn
	refreshClient *http.Client
	sendChan      chan []byte
	AuthToken     *oauth2.Token
	stop          chan bool
	pongRx        chan bool
	reconnectTime int
	isConnected bool
	mu *sync.Mutex
	channelPointHandlers map[string]func(*ChannelPointsEvent)
}

// Returns a new PubSubClient. 
func NewPubSubClient(config *oauth2.Config, userToken *oauth2.Token) *PubSubClient {
	return &PubSubClient{
		conn:          nil,
		AuthToken:     userToken,
		refreshClient: config.Client(context.Background(), userToken),
		sendChan:      make(chan []byte, 256),
		stop:          make(chan bool),
		pongRx:          make(chan bool, 1),
		reconnectTime: 1,
		isConnected: false,
		mu: &sync.Mutex{},
		channelPointHandlers: make(map[string]func(*ChannelPointsEvent)),
	}
}

// Thread-safe check of whether or not the PubSubClient is connected.
func (c *PubSubClient) IsConnected() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	result := c.isConnected
	return result
}

// Connect to the Twitch PubSub endpoint and listen on all registered topics. 
// Will automatically reconnect on failure.
// with exponential backoff. Returns an error if the client is already connected.
func (c *PubSubClient) Connect() error {
	if !c.IsConnected() {
		conn, _, err := websocket.DefaultDialer.Dial(pubSubURL, nil)
		if err != nil {
			go c.reconnect()
			return err
		}
		c.conn = conn
		go c.reader()
		go c.writer()

		// Listen on all registered topics
		var topics []string
		for id := range c.channelPointHandlers {
			topics = append(topics, channelPointTopic + id)
		}
		if len(topics) > 0 {
			c.listen(&topics)
		}

		c.mu.Lock()
		c.isConnected = true
		c.mu.Unlock()
		fmt.Println("PubSub Client connected")
	} else {
		errors.New("PubSub Client is already connected")
	}
	return nil
}

// Close disconnects the client from the Twitch PubSub endpoint.
// If the client is already connected, Close() will return an error.
func (c *PubSubClient) Close() error {
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
func (c *PubSubClient) shutdown() {
	c.Close()
	go c.reconnect()
}

// Updates the exponential backoff reconnect time and attempts to reconnect
// to the Twitch PubSub endpoint after this time period.
func (c *PubSubClient) reconnect() {
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
func (c *PubSubClient) reader() {
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
func (c *PubSubClient) handle(msg []byte) {
	builtMsg := new(PubSubMessage)
	err := json.Unmarshal([]byte(msg), builtMsg)
	if err != nil {
		fmt.Println(msg, err)
		return
	}

	switch builtMsg.Type {
	case "MESSAGE":
		split := strings.Split(builtMsg.Data.Topic, ".")
		topic := split[0] + "."
		id := split[1]
		switch topic {
		case channelPointTopic:
			event := new(ChannelPointsEvent)
			err = json.Unmarshal([]byte(builtMsg.Data.Message), event)
			if err != nil {
				fmt.Println(err)
				return
			}
			c.channelPointHandlers[id](event)
		default:
			fmt.Println("Unknown topic:", topic)
		}
	case "RESPONSE":
		resp := new(pubSubResponse)
		err := json.Unmarshal([]byte(msg), builtMsg)
		if err != nil {
			fmt.Println(msg, err)
			return
		}
		if resp.Error != "" {
			fmt.Println("PubSub client received error response:", resp.Error)
		}
	case "PONG":
		c.pongRx<- true
	case "RECONNECT":
		fmt.Println("PubSub client received reconnect message.")
	default:
		fmt.Println("Unknown message:", builtMsg.Type)
	} 
}

// writer handles transmitting regular ping messages and determines if a pong response is in time.
// Also writes any messages from the send channel.
func (c *PubSubClient) writer() {
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
func (c *PubSubClient) write(msg []byte) error {
	err := c.conn.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		fmt.Println("PubSub error in tx:", err)
		return err
	}
	return nil
}

// listen creates a listen request and sends it to the send channel.
func (c *PubSubClient) listen(topics *[]string) {
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
