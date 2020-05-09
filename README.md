# go-twitch-api

This library provides access to the [Helix Twitch API](https://dev.twitch.tv/docs/api/reference).

It handles both app access as well as user access tokens. All tokens used are automatically refreshed.

Note: This is a work in progress and a project to help me learn Go :). May not provide full functionality.

## Installation

Run:

```bash
$ go get github.com/kelr/go-twitch-api/twitchapi
```

## Usage
Example usage that gets the top active streamers:

```go
package main

import (
	"github.com/kelr/go-twitch-api/twitchapi"
	"encoding/json"
	"fmt"
)

// Provide your Client ID and secret here.
// Better to set these as environment variables.
const (
	clientID     = ""
	clientSecret = ""
)

func main() {

	client, err := twitchapi.NewTwitchClient(clientID, clientSecret)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	// Set options, English and only return the top 2 streams
	opt := &twitchapi.GetStreamsOpt{
		Language: "en",
		First:    2,
	}

	// Returns a GetStreamsResponse object
	response, err := client.GetStreams(opt)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Pretty print
	obj, _ := json.MarshalIndent(response, "", "  ")
	fmt.Println(string(obj))
}

```
