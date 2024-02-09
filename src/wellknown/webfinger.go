package wellknown

import (
	"net/url"
)

func GetWebfinger() Webfinger {
	logger.Trace("Getting webfinger")

	return Webfinger{
		
	}
}

type Webfinger struct{
	Subject string
	Links []Link
}

type Link struct {
	Rel string
	Type string
	Href url.URL
}