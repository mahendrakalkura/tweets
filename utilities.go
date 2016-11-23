package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"math/rand"
	"regexp"
	"sort"
	"strconv"
	"time"
)

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

func get_track(programs []Program) []string {
	var track []string

	var re *regexp.Regexp
	var matches []string

	re = regexp.MustCompile("\\w+")

	for _, program := range programs {
		matches = re.FindAllString(program.Query, -1)
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

func get_unix(value string) time.Time {
	integer, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		panic(err)
	}
	unix := time.Unix(integer, 0)
	return unix
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
