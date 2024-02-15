package activityPub

import "time"

type Activity struct {
	Id        string    `json:"id"`
	Type      string    `json:"type"`
	Context   any       `json:"context"`
	To        []string  `json:"to"`
	Bto       []string  `json:"bto"`
	Cc        []string  `json:"cc"`
	Bcc       []string  `json:"bcc"`
	Audience  []string  `json:"audience"`
	Object    any       `json:"object"`
	Published *time.Time `json:"published"`
	Actor     string    `json:"actor"`
}
