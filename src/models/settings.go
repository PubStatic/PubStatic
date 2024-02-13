package models

type Settings struct {
	ActivityPubSettings ActivityPubSettings `yaml:"activityPub"`
	ServerSettings      ServerSettings      `yaml:"server"`
}

type ServerSettings struct {
	Port int `yaml:"port"`
}

type ActivityPubSettings struct {
	UserName         string `yaml:"userName"`
	CosmeticUserName string `yaml:"cosmeticName"`
	UserDescription  string `yaml:"userDescription"`
}
