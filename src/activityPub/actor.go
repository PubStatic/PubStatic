package activityPub

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetActor(host string, preferredUsername string, name string, summary string, publicKeyPem string) Actor {
	logger.Trace("Getting Actor")

	id := fmt.Sprintf("https://%s/%s", host, preferredUsername)

	return Actor{
		Context: []string{
			"https://www.w3.org/ns/activitystreams",
		},
		Type:              "Person",
		Id:                id,
		Outbox:            id + "/outbox",
		Inbox:             id + "/inbox",
		Following:         id + "/following",
		Followers:         id + "/followers",
		PreferredUsername: preferredUsername,
		Name:              name,
		Summary:           summary,
		Icon: []string{
			"", // TODO Add icon url here
		},
		PublicKey: PublicKey{
			Context:      "https://w3id.org/security/v1",
			Type:         "Key",
			Id:           id + "#main-key",
			Owner:        id,
			PublicKeyPem: publicKeyPem,
		},
	}
}

type Actor struct {
	Context           []string `json:"@context"`
	Type              string
	Id                string
	Outbox            string
	Following         string
	Followers         string
	Inbox             string
	PreferredUsername string
	Name              string
	Summary           string
	Icon              []string
	PublicKey         PublicKey
}

type PublicKey struct {
	Context      string `json:"@context"`
	Type         string `json:"@type"`
	Id           string
	Owner        string
	PublicKeyPem string
}

func GetForeignActor(actorId string) (*Actor, error) {
	logger.Tracef("Getting foreign actor with id: %s", actorId)

	// Send an HTTP GET request
	response, err := http.Get(actorId)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Read the response body
	body, err := io.ReadAll(io.Reader(response.Body))
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON response into an Actor struct
	var actor Actor
	err = json.Unmarshal(body, &actor)
	if err != nil {
		return nil, err
	}

	return &actor, nil
}
