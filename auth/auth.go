package auth

import (
	"encoding/json"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/twitch"
	"io"
	"os"
)

const (
	stateLen = 32
)

// Creates and returns OAuth2 configuration object with the twitch endpoint. Also returns a URL to be sent to the user used to initiate authentication.
func NewUserAuth(clientID string, clientSecret string, redirectURI string, scopes *[]string) (*oauth2.Config, error) {
	if clientID == "" {
		return nil, errors.New("A Client ID must be provided to create an OAuth2 config")
	}
	if clientSecret == "" {
		return nil, errors.New("A Client secret must be provided to create a OAuth2 config")
	}
	if redirectURI == "" {
		return nil, errors.New("A redirect URI must be provided to create a OAuth2 config")
	}
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       *scopes,
		Endpoint:     twitch.Endpoint,
		RedirectURL:  redirectURI,
	}
	return config, nil
}

// Returns a URL to send to the end user for them to access as well as the state string embedded into the URL. Ensure that this state string matches the value recieved at the redirect URI.
func GetAuthCodeURL(config *oauth2.Config) (string, string) {
	state, _ := generateState()
	return config.AuthCodeURL(state, oauth2.AccessTypeOffline), state
}

// Conducts the exchange to turn an auth code into a user token. The OAuth2 config used to create the auth code must be the same.
func TokenExchange(config *oauth2.Config, authCode string) (*oauth2.Token, error) {
	token, err := config.Exchange(context.Background(), authCode)
	if err != nil {
		fmt.Println("Error in obtaining user token:", err)
		return nil, err
	}
	return token, nil
}

func FlushTokenFile(tokenFile string, token *oauth2.Token) error {
	f, err := os.OpenFile(tokenFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer f.Close()
	return encodeFile(f, token)
}

func LoadTokenFile(config *oauth2.Config, tokenFile string) (*oauth2.Token, error) {
	f, err := os.OpenFile(tokenFile, os.O_RDWR, 0755)
	if err != nil {
		// Handle PathError specifically as it indicates the file does not exist
		if _, ok := err.(*os.PathError); ok {
			fmt.Println("No saved token file found.")
			os.Exit(0)
		} else {
			return nil, err
		}
	}
	defer f.Close()
	token, err := decodeFile(f)
	if err != nil {
		return nil, err
	}

	return token, nil
}

// VerifyToken checks the input token for validity, and will return a refreshed token
// if it has expired. If the token is still valid, will return the same token.
func VerifyToken(config *oauth2.Config, oldToken *oauth2.Token) *oauth2.Token {
	tokenSource := config.TokenSource(oauth2.NoContext, oldToken)
	newToken, err := tokenSource.Token()
	if err != nil {
	    fmt.Println(err)
	}
	return newToken
}

func encodeFile(file *os.File, token *oauth2.Token) error {
	enc := json.NewEncoder(file)
	enc.SetIndent("", " ")
	err := enc.Encode(token)
	if err == io.EOF {
		err = nil
	}
	return err
}

func decodeFile(file *os.File) (*oauth2.Token, error) {
	token := new(oauth2.Token)
	err := json.NewDecoder(file).Decode(token)
	if err == io.EOF {
		err = nil
	}
	return token, err
}

// Generate random 32 character state string
func generateState() (string, error) {
	var buf [stateLen]byte
	if _, err := rand.Read(buf[:]); err != nil {
		return "", err
	}

	return hex.EncodeToString(buf[:]), nil
}
