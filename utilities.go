package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"time"
)

func get_database(settings *Settings) *sqlx.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		settings.SQLX.Hostname,
		settings.SQLX.Port,
		settings.SQLX.Username,
		settings.SQLX.Password,
		settings.SQLX.Database,
	)
	database := sqlx.MustConnect("postgres", dsn)
	return database
}

func get_http_client(settings *Settings, with_proxy bool) *http.Client {
	timeout := time.Duration(30 * time.Second)

	proxy := get_proxy(settings.Proxies.Hostname, settings.Proxies.Ports)
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

func get_number(value string) string {
	re := regexp.MustCompile("[^0-9]")
	value = re.ReplaceAllString(value, "")
	return value
}

func get_proxy(hostname string, ports []int) string {
	port := get_random_number(ports[0], ports[1]+1)
	return fmt.Sprintf("https://%s:%d", hostname, port)
}

func get_random_number(minimum int, maximum int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(maximum-minimum) + minimum
}

func get_settings() *Settings {
	var settings = &Settings{}
	_, err := toml.DecodeFile("settings.toml", settings)
	if err != nil {
		panic(err)
	}
	return settings
}

func get_timestamp_from_integer(value string) time.Time {
	integer, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		panic(err)
	}
	timestamp := time.Unix(integer, 0)
	return timestamp
}

func get_timestamp_from_string_1(value string) time.Time {
	timestamp, err := time.Parse(time.RubyDate, value)
	if err != nil {
		panic(err)
	}
	return timestamp
}

func get_timestamp_from_string_2(value string) time.Time {
	var timestamp time.Time
	var err error

	timestamp, err = time.Parse("3:04 PM - 2 Jan 2006", value)
	if err == nil {
		return timestamp
	}

	timestamp, err = time.Parse("15:04 - 2. Jan. 2006", value)
	if err == nil {
		return timestamp
	}

	return time.Now().UTC()
}

func get_track(programs []Program) []string {
	var track []string
	re := regexp.MustCompile("\\w+")
	for _, program := range programs {
		matches := re.FindAllString(program.QueriesTwitter, -1)
		for _, match := range matches {
			if len(match) < 3 {
				continue
			}
			if len(match) > 25 {
				continue
			}
			if in_array(match, track) {
				continue
			}
			track = append(track, match)
		}
	}
	sort.Sort(ByLengthAndValue(track))
	if len(track) > 400 {
		track = track[0:400]
	}
	return track
}

func get_twitter_client(
	consumer_key string, consumer_secret string, access_key string, access_secret string,
) *twitter.Client {
	oauth1_config := oauth1.NewConfig(consumer_key, consumer_secret)
	oauth1_token := oauth1.NewToken(access_key, access_secret)
	oauth1_client := oauth1_config.Client(oauth1.NoContext, oauth1_token)

	client := twitter.NewClient(oauth1_client)

	return client
}

func get_twitter_parameters(track []string) *twitter.StreamFilterParams {
	parameters := &twitter.StreamFilterParams{
		Track: track,
	}
	return parameters
}

func has_stopped(program *Program) bool {
	now := time.Now().UTC()
	if now.Before(program.BeginningAt) {
		return true
	}
	if now.After(program.EndingAt) {
		return true
	}
	return false
}

func in_array(value string, array []string) bool {
	for key := range array {
		ok := array[key] == value
		if ok {
			return true
		}
	}
	return false
}
