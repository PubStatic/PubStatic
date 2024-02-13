package activityPub

import (
	"errors"

	"github.com/PubStatic/PubStatic/repository"
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
		follow(activity, connectionString)
	case "Undo":
		undo(activity)
	}

	return nil
}

func follow(activity Activity, connectionString string) error {
	logger.Debug(activity)

	return repository.WriteMongo("Inbox", "Follow", activity, connectionString)
}

func undo(activity Activity) {

}
