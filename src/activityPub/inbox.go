package activityPub

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/PubStatic/PubStatic/repository"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
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
		return follow(activity, connectionString, *actor, host)
	case "Undo":
		return undo(activity, connectionString, *actor)
	}

	return nil
}

func follow(activity Activity, connectionString string, actor Actor, ownHost string) error {
	count, err := repository.CountMongo[Activity]("Inbox", "Follow", bson.D{{"actor", actor.Id}}, connectionString)

	if err != nil {
		return err
	}

	if count != 0 {
		return fmt.Errorf("follow for this actor already exists. actorId: %s", actor.Id)
	}

	url, err := url.Parse(actor.Inbox)

	if err != nil {
		return err
	}

	time := time.Now()

	err = SendActivity(Activity{
		Context:   "https://www.w3.org/ns/activitystreams",
		Id:        fmt.Sprintf("https://%s/accept/%s", ownHost, uuid.New()),
		Type:      "Accept",
		Actor:     fmt.Sprintf("https://%s", ownHost),
		Object:    activity,
		Published: &time,
	}, *url, ownHost, connectionString)

	if err != nil {
		logger.Warn("Sending activity failed")

		return err
	}

	return repository.WriteMongo("Inbox", "Follow", activity, connectionString)
}

func undo(activity Activity, connectionString string, actor Actor) error {
	logger.Trace("Entered undo")

	// TODO Check if it is really a undo follow
	
	count, err := repository.CountMongo[Activity]("Inbox", "Follow", bson.D{{"actor", actor.Id}}, connectionString)

	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("cannot undo follow. follow does not exist. actorId: %s", actor.Id)
	}

	deleteCount, err := repository.DeleteMongo("Inbox", "Follow", bson.D{{"actor", actor.Id}}, connectionString)

	if err != nil {
		return err
	}

	logger.Debugf("Deleted: %d items", deleteCount)

	return nil
}
