package activityPub

import (
	"fmt"
	"net/http"
)

func sendActivity(activity Activity, inbox string) error {

	logger.Trace("Sending activity")

	req, err := http.NewRequest("POST", inbox, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Accept", "application/json")

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return fmt.Errorf("sending activity failed with error code: %d", response.StatusCode)
	}

	logger.Trace("Successfully send activity")

	return nil
}
