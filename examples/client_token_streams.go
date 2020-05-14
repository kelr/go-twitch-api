// An example to create a client with a client credentials token.
// Uses the token to get info about a stream.
package main

import (
	"../helix"
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
