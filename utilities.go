package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func get_client(settings *Settings, with_proxy bool) *http.Client {
	timeout := time.Duration(30 * time.Second)

	proxy := get_proxy(settings)
	proxy_url, err := url.Parse(proxy)
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	client.Timeout = timeout
	if with_proxy {
		client.Transport = &http.Transport{Proxy: http.ProxyURL(proxy_url)}
	}

	return client
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

func get_settings() *Settings {
	var settings = &Settings{}
	_, err := toml.DecodeFile("settings.toml", settings)
	if err != nil {
		panic(err)
	}
	return settings
}

func get_proxy(settings *Settings) string {
	port := get_random(settings.Proxies.Ports[0], settings.Proxies.Ports[1]+1)
	return fmt.Sprintf("https://%s:%d", settings.Proxies.Hostname, port)
}

func get_random(minimum int, maximum int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(maximum-minimum) + minimum
}

func get_unix(value string) time.Time {
	integer, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		panic(err)
	}
	unix := time.Unix(integer, 0)
	return unix
}
