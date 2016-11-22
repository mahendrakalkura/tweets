package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Settings struct {
	SQLX    SettingsSQLX    `toml:"sqlx"`
	Twitter SettingsTwitter `toml:"twitter"`
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

func get_settings() *Settings {
	var settings = &Settings{}
	_, err := toml.DecodeFile("settings.toml", settings)
	if err != nil {
		panic(err)
	}
	return settings
}

func get_database(settings_sqlx *SettingsSQLX) *sqlx.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		settings_sqlx.Hostname,
		settings_sqlx.Port,
		settings_sqlx.Username,
		settings_sqlx.Password,
		settings_sqlx.Database,
	)
	database := sqlx.MustConnect("postgres", dsn)
	return database
}
