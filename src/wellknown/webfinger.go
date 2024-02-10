package wellknown

import (
	"fmt"
)

func GetWebfinger(host string, userName string) Webfinger {
	logger.Trace("Getting webfinger")

	logger.Debugf("Host: %s", host)
	logger.Debugf("UserName: %s", userName)

	return Webfinger{
		Subject: fmt.Sprintf("acct:%s@%s", userName, host),
		Links: []Link{
			{
				Rel:  "self",
				Type: "application/activity+json",
				Href: fmt.Sprintf("https://%s/%s", host, userName),
			},
		},
	}
}

type Webfinger struct {
	Subject string `json:"subject"`
	Links   []Link `json:"links"`
}

type Link struct {
	Rel  string `json:"rel"`
	Type string `json:"type"`
	Href string `json:"href"`
}
