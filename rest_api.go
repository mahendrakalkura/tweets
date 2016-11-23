package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func rest_api(settings *Settings, database *sqlx.DB) {
	fmt.Println("rest_api() - Start")

	channels_program_and_max_position := make(chan ProgramAndMaxPosition)
	channels_screen_name := make(chan string)
	channels_signal := make(chan os.Signal)
	channels_exit := make(chan bool)

	for index := 1; index <= 64; index++ {
		go rest_api_consumer_tweets(settings, database, channels_program_and_max_position, channels_screen_name)
		go rest_api_consumer_tweeters(settings, database, channels_screen_name)
	}

	go rest_api_producer(database, channels_program_and_max_position, channels_screen_name)

	ticker := time.NewTicker(60 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				go rest_api_producer(database, channels_program_and_max_position, channels_screen_name)
			}
		}
	}()

	signal.Notify(channels_signal, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)

	<-channels_signal

	close(channels_program_and_max_position)
	close(channels_screen_name)

	channels_exit <- true

	<-channels_exit

	fmt.Println("rest_api() - Stop")
}

func rest_api_consumer_tweets(
	settings *Settings,
	database *sqlx.DB,
	channels_program_and_max_position chan ProgramAndMaxPosition,
	channels_screen_name chan string,
) {
	fmt.Println("rest_api_consumer_tweets() - Start")

	for program_and_max_position := range channels_program_and_max_position {
		if has_stopped(program_and_max_position.Program) {
			continue
		}
		tweets, max_position := get_tweets(
			settings, program_and_max_position.Program.Query, program_and_max_position.MaxPosition,
		)
		for _, tweet := range tweets {
			go rest_api_set_tweet(database, tweet, channels_screen_name)
		}
		if max_position != "" {
			var program_and_max_position = ProgramAndMaxPosition{
				Program:     program_and_max_position.Program,
				MaxPosition: max_position,
			}
			channels_program_and_max_position <- program_and_max_position
		}
	}

	fmt.Println("rest_api_consumer_tweets() - Stop")
}

func rest_api_consumer_tweeters(settings *Settings, database *sqlx.DB, channels_screen_name chan string) {
	fmt.Println("rest_api_consumer_tweeters() - Start")

	for screen_name := range channels_screen_name {
		tweeter := get_tweeters(settings, screen_name)
		go set_tweeter(database, &tweeter)
	}

	fmt.Println("rest_api_consumer_tweeters() - Stop")
}

func rest_api_producer(
	database *sqlx.DB,
	channels_program_and_max_position chan ProgramAndMaxPosition,
	channels_screen_name chan string,
) {
	fmt.Println("rest_api_producer() - Start")

	programs := get_programs(database)
	for _, program := range programs {
		var program_and_max_position = ProgramAndMaxPosition{
			Program:     program,
			MaxPosition: "",
		}
		channels_program_and_max_position <- program_and_max_position
	}

	screen_names := get_screen_names(database)
	for _, screen_name := range screen_names {
		channels_screen_name <- screen_name
	}

	fmt.Println("rest_api_producer() - Stop")
}

func rest_api_set_tweet(database *sqlx.DB, tweet Tweet, channels_screen_name chan string) {
	set_tweet(database, &tweet)
	tweeter, err := get_tweeter(database, tweet.UserScreenName)
	if err == nil {
		set_tweeter(database, tweeter)
	} else {
		channels_screen_name <- tweet.UserScreenName
	}
}
