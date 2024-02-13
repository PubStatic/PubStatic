package activityPub

import (
	"errors"
	"fmt"
	"github.com/PubStatic/PubStatic/repository"
	"github.com/google/uuid"
)

func ReceiveActivity(activity Activity, header map[string][]string, host string, connectionString string) error {

	logger.Trace("Entered ReceiveActivity")

	actor, err := GetForeignActor(activity.Actor)
	if err != nil {
		logger.Warn("Could not get foreign actor")
		return err
	}

	isSignatureValid, err := validateSignature(header, actor.PublicKey, host)

	if err != nil || !isSignatureValid {
		return errors.New("signature validation failed")
	}

	switch activity.Type {
	case "Follow":
		follow(activity, connectionString, *actor, host)
	case "Undo":
		undo(activity)
	}

	return nil
}

func follow(activity Activity, connectionString string, actor Actor, host string) error {

	sendActivity(Activity{
		Context: "https://www.w3.org/ns/activitystreams",
		Id:      fmt.Sprintf("https://%s/accept/%s", host, uuid.New()),
		Type:    "Accept",
		Actor:   fmt.Sprintf("https://%s", host),
		Object:  activity,
	}, actor.Inbox)

	return repository.WriteMongo("Inbox", "Follow", activity, connectionString)
}

func undo(activity Activity) {

}
