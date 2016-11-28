package main

import (
	"time"
)

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

type ItemsAndMaxPosition struct {
	Items       string `json:"items_html"`
	MaxPosition string `json:"min_position"`
}

type Program struct {
	Query       string    `db:"query"`
	BeginningAt time.Time `db:"beginning_at"`
	EndingAt    time.Time `db:"ending_at"`
}

type ProgramAndMaxPosition struct {
	Program     Program
	MaxPosition string
}

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

type Tweet struct {
	Id                  string     `db:"twitter_id"`
	Text                string     `db:"twitter_text"`
	Retweets            int        `db:"twitter_retweets"`
	Timestamp           time.Time  `db:"twitter_timestamp"`
	UserId              string     `db:"twitter_user_id"`
	UserScreenName      string     `db:"twitter_user_screen_name"`
	UserName            string     `db:"twitter_user_name"`
	UserProfileImageURL string     `db:"twitter_user_profile_image_url"`
	UserTweets          *int       `db:"twitter_user_tweets"`
	UserFollowers       *int       `db:"twitter_user_followers"`
	UserFollowing       *int       `db:"twitter_user_following"`
	UserTimestamp       *time.Time `db:"twitter_user_timestamp"`
}

type Tweeter struct {
	ScreenName string    `db:"twitter_user_screen_name"`
	Tweets     int       `db:"twitter_user_tweets"`
	Followers  int       `db:"twitter_user_followers"`
	Following  int       `db:"twitter_user_following"`
	Timestamp  time.Time `db:"twitter_user_timestamp"`
}
