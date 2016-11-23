package main

import (
	"time"
)

type Settings struct {
	Proxies SettingsProxies `toml:"proxies"`
	SQLX    SettingsSQLX    `toml:"sqlx"`
	Twitter SettingsTwitter `toml:"twitter"`
}

type SettingsProxies struct {
	Hostname string `toml:"hostname"`
	Ports    []int  `toml:"ports"`
}

type SettingsSQLX struct {
	Database string `toml:"database"`
	Hostname string `toml:"hostname"`
	Password string `toml:"password"`
	Port     string `toml:"port"`
	Username string `toml:"username"`
}

type SettingsTwitter struct {
	ConsumerKey    string `toml:"consumer_key"`
	ConsumerSecret string `toml:"consumer_secret"`
	AccessKey      string `toml:"access_key"`
	AccessSecret   string `toml:"access_secret"`
}

type Program struct {
	Query       string    `db:"query"`
	BeginningAt time.Time `db:"beginning_at"`
	EndingAt    time.Time `db:"ending_at"`
}

type Body struct {
	Items       string `json:"items_html"`
	MaxPosition string `json:"min_position"`
}

type Tweet struct {
	Id                  string
	Text                string
	Retweets            int
	CreatedAt           time.Time
	UserId              string
	UserScreenName      string
	UserName            string
	UserProfileImageUrl string
	UserTweets          int
	UserFollowers       int
	UserFollowing       int
	UserCreatedAt       time.Time
}

type ByLengthAndValue []string

func (items ByLengthAndValue) Len() int {
	return len(items)
}

func (items ByLengthAndValue) Swap(one int, two int) {
	items[one], items[two] = items[two], items[one]
}

func (items ByLengthAndValue) Less(one, two int) bool {
	if len(items[one]) < len(items[two]) {
		return true
	}
	if len(items[one]) > len(items[two]) {
		return false
	}
	return items[one] < items[two]
}
