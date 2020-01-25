# go-twitch-api

This library provides access to the Helix Twitch API.

Note: This is a work in progress and may not provide full functionality.

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
	"./twitchapi"
	"encoding/json"
	"fmt"
)

// Provide your Client ID here
const clientID = ""

func main() {

	client := twitchapi.NewTwitchClient(clientID)

	// Set options, English and only return the top 2 streams
	opt := &twitchapi.GetStreamsOpt{
		Language: "en",
		First:    2,
	}

	// Returns a GetStreamsResponse object
	response, err := client.GetStreams(opt)
	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Println(response.Data, response.Pagination.Cursor)

	// Pretty print
	obj, _ := json.MarshalIndent(response, "", "  ")
	fmt.Println(string(obj))
}


```
