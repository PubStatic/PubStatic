package wellknown

func GetNodeInfo2_1(version string) NodeInfo {
	logger.Trace("Getting NodeInfo")

	return NodeInfo{
		Version: "2.1",
		Software: Software{
			Name:       "PubStatic",
			Version:    version,
			Repository: "https://github.com/PubStatic/PubStatic",
			Homepage:   "https://lna-dev.net",
		},
		Protocols: []string{
			"activitypub",
		},
		Services: Service{
			Inbound:  []Service{},
			Outbound: []Service{},
		},
		Usage: Usage{
			LocalPosts:    0, // TODO this need to be filled
			LocalComments: 0, // TODO this need to be filled
			Users: Users{
				ActiveHalfyear: 1,
				ActiveMonth:    1,
				Total:          1,
			},
		},
		OpenRegistration: false,
		Metadata:         map[string]string{},
	}
}

type NodeInfo struct {
	Version          string            `json:"version"`
	Software         Software          `json:"software"`
	Protocols        []string          `json:"protocols"`
	Services         Service           `json:"services"`
	Usage            Usage             `json:"usage"`
	OpenRegistration bool              `json:"openRegistration"`
	Metadata         map[string]string `json:"metadata"`
}

type Software struct {
	Name       string `json:"name"`
	Version    string `json:"version"`
	Repository string `json:"repository"`
	Homepage   string `json:"homepage"`
}

type Service struct {
	Outbound []Service `json:"outbound"` // Empty array
	Inbound  []Service `json:"inbound"` // Empty array
}

type Usage struct {
	LocalPosts    int `json:"localPosts"`
	LocalComments int `json:"localComments"`
	Users         Users `json:"users"`
}

type Users struct {
	ActiveHalfyear int `json:"activeHalfyear"`
	ActiveMonth    int `json:"activeMonth"`
	Total          int `json:"total"`
}
