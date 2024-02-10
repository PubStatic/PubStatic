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
	Version          string
	Software         Software
	Protocols        []string
	Services         Service
	Usage            Usage
	OpenRegistration bool
	Metadata         map[string]string
}

type Software struct {
	Name       string
	Version    string
	Repository string
	Homepage   string
}

type Service struct {
	Outbound []Service // Empty array
	Inbound  []Service // Empty array
}

type Usage struct {
	LocalPosts    int
	LocalComments int
	Users         Users
}

type Users struct {
	ActiveHalfyear int
	ActiveMonth    int
	Total          int
}
