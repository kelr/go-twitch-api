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
	"fmt"
	"github.com/kelr/go-twitch-api/twitchapi"
)

// Provide your Client ID here
const clientID = ""

func main() {

	client := twitchapi.NewTwitchClient(clientID)

	response, err := client.GetStreams(nil)
	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Printf("%v\n%s\n", response.Data[0], response.Pagination.Cursor)
}

```
