package activityPub

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetActor(host string, preferredUsername string, name string, summary string, publicKeyPem string) Actor {
	logger.Trace("Getting Actor")

	id := fmt.Sprintf("https://%s", host)

	return Actor{
		Context: []string{
			"https://www.w3.org/ns/activitystreams",
			"https://w3id.org/security/v1",
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
		Icon:              []string{}, // TODO Add icon url here
		PublicKey: PublicKey{
			Id:           id + "#main-key",
			Owner:        id,
			PublicKeyPem: publicKeyPem,
		},
	}
}

type Actor struct {
	Context           []string  `json:"@context"`
	Type              string    `json:"type"`
	Id                string    `json:"id"`
	Outbox            string    `json:"outbox"`
	Following         string    `json:"following"`
	Followers         string    `json:"followers"`
	Inbox             string    `json:"inbox"`
	PreferredUsername string    `json:"preferredUsername"`
	Name              string    `json:"name"`
	Summary           string    `json:"summary"`
	Icon              []string  `json:"icon"`
	PublicKey         PublicKey `json:"publicKey"`
}

type PublicKey struct {
	Id           string `json:"id"`
	Owner        string `json:"owner"`
	PublicKeyPem string `json:"publicKeyPem"`
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
