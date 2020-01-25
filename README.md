# go-twitch-api

This library provides access to the Helix Twitch API.

## Installation

Run:

```bash
$ go get github.com/kelr/go-twitch-api/twitchapi
```


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


## License

All files under this repository fall under the MIT License (see the file LICENSE).