package pubsub

import (
	"fmt"
	"context"
	"encoding/json"
	"github.com/gorilla/websocket"
	"golang.org/x/oauth2"
	"math/rand"
	"net/http"
	"strings"
)

const (
	pubSubURL = "wss://pubsub-edge.twitch.tv"
	nonceSet  = "abcdedfghijklmnopqrstABCDEFGHIJKLMNOP123456789"
)

type PubSubClient struct {
	conn          *websocket.Conn
	IsConnected   bool
	refreshClient *http.Client
	messages      chan string
	AuthToken     *oauth2.Token
}

func NewPubSubClient(config *oauth2.Config, userToken *oauth2.Token) (*PubSubClient) {
	return &PubSubClient{
		conn:          nil,
		IsConnected:   false,
		AuthToken:     userToken,
		refreshClient: config.Client(context.Background(), userToken),
		Messages:      make(chan<- string, 1),
	}
}

func (c *PubSubClient) Connect() error {
	conn, _, err := websocket.DefaultDialer.Dial(pubSubURL, nil)
	if err != nil {
		fmt.Println("PubSub error in dial:", err)
		return err
	}
	c.conn = conn

	fmt.Println("PubSub connected")
	c.IsConnected = true

	go func() {
		defer close(c.messages)
		for {
			_, msg, err := c.conn.ReadMessage()
			if err != nil {
				fmt.Println("PubSub error in rx:", err)
				c.IsConnected = false
				return
			}
			ret := string(msg)
			fmt.Println("RX:", ret)
			c.messages <- ret
		}
	}()
	return nil
}

func (c *PubSubClient) write(msg []byte) {
	err := c.conn.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		fmt.Println("PubSub error in tx:", err)
		c.IsConnected = false
		return
	}
}

type pubSubResponse struct {
	Type  string `json:"type"`
	Nonce string `json:"nonce"`
	Error string `json:"error"`
}

type pubSubMessageResponse struct {
	Type  string `json:"type"`
	Data  struct {
		Topic    string `json:"topic"`
		Message string   `json:"message"`
	} `json:"data"`
}

type pubSubListenRequest struct {
	Type  string `json:"type"`
	Nonce string `json:"nonce"`
	Data  struct {
		Topics    []string `json:"topics"`
		AuthToken string   `json:"auth_token"`
	} `json:"data"`
}

func (c *PubSubClient) Listen(topics *[]string) {
	request := pubSubListenRequest{
		Type:  "LISTEN",
		Nonce: generateNonce(15),
	}
	request.Data.Topics = *topics
	request.Data.AuthToken = c.AuthToken.AccessToken
	bytes, _ := json.Marshal(request)
	c.write(bytes)
}

func generateNonce(length int) string {
	var curr strings.Builder
	for i := 0; i < length; i++ {
		curr.WriteString(string(nonceSet[rand.Intn(len(nonceSet))]))
	}
	return curr.String()
}
