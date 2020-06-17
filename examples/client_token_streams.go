// An example to create a client with the client credentials OAuth flow.
// Gets info about streams.
package main

import (
	"encoding/json"
	"log"
	"github.com/kelr/gundyr/helix"
)

// Provide your Client ID and secret here.
// Better to set these as environment variables.
const (
	clientID     = ""
	clientSecret = ""
)

func main() {
	client, err := helix.NewClient(clientID, clientSecret)
	if err != nil {
		log.Fatal(err)
	}

	// Get the top 2 english streams. Response is a GetStreamsResponse object.
	response, err := client.GetStreams(&helix.GetStreamsOpt{Language: "en", First: 2})
	if err != nil {
		log.Fatal(err)
	}

	// Pretty print
	obj, _ := json.MarshalIndent(response, "", "  ")
	log.Println(string(obj))
}
