package wellknown

import "fmt"

func GetLinkToNodeInfo(host string) NodeInfoLink{
	logger.Trace("Getting link to NodeInfo")

	return NodeInfoLink{
		Links: []NodeInfoSubLink{
			{
				Rel: "http://nodeinfo.diaspora.software/ns/schema/2.1",
				Href: fmt.Sprintf("https://%s/nodeinfo/2.1", host),
			},
		},
	}
}

type NodeInfoLink struct{
	Links []NodeInfoSubLink `json:"name"`
}

type NodeInfoSubLink struct{
	Rel string `json:"rel"`
	Href string `json:"href"`
}
