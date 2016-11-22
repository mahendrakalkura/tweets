package main

import (
    "time"
)

type Settings struct {
    Proxies SettingsProxies `toml:"proxies"`
    SQLX    SettingsSQLX    `toml:"sqlx"`
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
