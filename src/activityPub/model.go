package activityPub

import "time"

type Activity struct {
	Id        string
	Type      string
	Context   string
	To        []string
	Bto       []string
	Cc        []string
	Bcc       []string
	Audience  []string
	Object    ActivityObject
	Published time.Time
	Actor     string
}

type ActivityObject interface {
}
