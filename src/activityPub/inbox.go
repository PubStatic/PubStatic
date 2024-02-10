package activityPub

import "errors"

func ReceiveActivity(activity Activity, header map[string][]string) error {

	actor, err := GetForeignActor(activity.Actor)
	if err != nil {
		return err
	}

	isSignatureValid, err := validateSignature(header, actor.PublicKey)

	if err != nil || !isSignatureValid {
		return errors.New("invalid signature")
	}

	switch activity.Type {
	case "Follow":
		follow(activity)
	case "Undo":
		undo(activity)
	}

	return nil
}

func follow(activity Activity) {

}

func undo(activity Activity) {

}
