# gundyr 

[![GoDoc](https://godoc.org/github.com/kelr/gundyr?status.png)](https://godoc.org/github.com/kelr/gundyr) [![Go Report Card](https://goreportcard.com/badge/github.com/kelr/gundyr)](https://goreportcard.com/report/github.com/kelr/gundyr)

Gundyr provides an easy to use interface to the [Helix Twitch API](https://dev.twitch.tv/docs/api/reference) and [Twitch PubSub](https://dev.twitch.tv/docs/pubsub).

It handles both app access as well as user access tokens. All tokens used are automatically refreshed.

Note: This is a work in progress and a project to help me learn Go :). May not provide full functionality.

## Install

```bash
$ go get github.com/kelr/gundyr
```

## Docs

Documentation can be found at [godoc](https://godoc.org/github.com/kelr/gundyr). Examples can be found in the [examples](https://github.com/kelr/gundyr/tree/master/examples) directory.

## Usage
### Helix - Getting a User ID

```go
cfg := &gundyr.HelixConfig{
	ClientID:     clientID, 
	ClientSecret: clientSecret,
}

c, err := gundyr.NewHelix(cfg)
if err != nil {
	log.Fatal(err)
}

userID, err := c.UserToID("kyrotobi")
if err != nil {
	log.Fatal(err)
}
log.Println(userID)
```

### Helix - Using User Access Tokens

If an OAuth2 token is not provided to HelixConfig, authentication will attempt to use the OAuth2 Client Credentials flow to obtain a App Access token.

```go
// See examples/auth_token.go on creating/retrieving tokens.
cfg := &gundyr.HelixConfig{
	ClientID:     clientID,
	ClientSecret: clientSecret,
	Scopes:       []string{"user:read:email"},
	RedirectURI:  redirectURI,
	Token:        token,
}

c, err := gundyr.NewHelix(cfg)
if err != nil {
	log.Fatal(err)
}

email, err := c.GetUserEmail("your-username")
if err != nil {
	log.Fatal(err)
}
log.Println(email)
```

### PubSub - Subscribing to Channel Point Redemptions

```go
func handleChannelPoints(event *pubsub.ChannelPointsEvent) {
	fmt.Println(event.Data.Redemption.Reward.Title)
}

func main() {
	scopes := []string{"channel:read:redemptions"}

	// Setup OAuth2 configuration
	config, err := auth.NewUserAuth(clientID, clientSecret, redirectURI, &scopes)
	if err != nil {
		log.Fatal(err)
	}

	// See examples/auth_token.go for an example on creating a new token.
	token, err := auth.RetrieveTokenFile(config, tokenFile)
	if err != nil {
		log.Fatal(err)
	}

	// Create a PubSub client and listen to the topics.
	client := pubsub.NewClient(config, token)
	client.ListenChannelPoints(userID, handleChannelPoints)
	client.Connect()
	select {}
}

```

## Contributions
Any and all contributions or bug fixes are appreciated.
