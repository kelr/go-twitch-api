# gundyr

This library provides an easy to use interface to the [Helix Twitch API](https://dev.twitch.tv/docs/api/reference).

It handles both app access as well as user access tokens. All tokens used are automatically refreshed.

Note: This is a work in progress and a project to help me learn Go :). May not provide full functionality.

## Installation

Run:

```bash
$ go get github.com/kelr/gundyr/
```

## Usage
Example to convert a Twitch login name to account ID

```go
package main

import (
	"fmt"
	"github.com/kelr/gundyr"
)

// Provide your Client ID and secret here.
const (
	clientID     = ""
	clientSecret = ""
)

func main() {
	client, err := helix.NewHelixClient(clientID, clientSecret)
	if err != nil {
		fmt.Println(err)
		return
	}

	userId, err := client.UserToId("kyrotobi")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(userId)
}
```

If more granular access is needed, import the helix package directly.
This example gets the top 2 live streamers.

```go
package main

import (
	"github.com/kelr/gundyr/helix"
	"encoding/json"
	"fmt"
)

// Provide your Client ID and secret here.
const (
	clientID     = ""
	clientSecret = ""
)

func main() {

	client, err := helix.NewTwitchClient(clientID, clientSecret)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	// Set options, English and only return the top 2 streams
	opt := &helix.GetStreamsOpt{
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
