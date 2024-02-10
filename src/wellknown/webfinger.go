package wellknown

import (
	"fmt"
)

func GetWebfinger(host string, userName string) Webfinger {
	logger.Trace("Getting webfinger")

	logger.Debug(fmt.Sprintf("Host: %s\nUserName: %s", host, userName))

	return Webfinger{
		Subject: fmt.Sprintf("acct:%s@%s", userName, host),
		Links: []Link{
			{
				Rel:  "self",
				Type: "application/activity+json",
				Href: fmt.Sprintf("https://%s", host),
			},
		},
	}
}

type Webfinger struct {
	Subject string
	Links   []Link
}

type Link struct {
	Rel  string
	Type string
	Href string
}
