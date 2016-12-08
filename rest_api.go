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

	channels_program_and_max_position := make(chan ProgramAndMaxPosition, 128)
	channels_signal := make(chan os.Signal)

	for index := 1; index <= 64; index++ {
		go rest_api_consumer(settings, database, channels_program_and_max_position)
	}

	go rest_api_producer(database, channels_program_and_max_position)

	ticker := time.NewTicker(60 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				go rest_api_producer(database, channels_program_and_max_position)
			}
		}
	}()

	signal.Notify(channels_signal, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)

	<-channels_signal

	close(channels_program_and_max_position)

	fmt.Println("rest_api() - Stop")
}

func rest_api_consumer(
	settings *Settings, database *sqlx.DB, channels_program_and_max_position chan ProgramAndMaxPosition,
) {
	fmt.Println("rest_api_consumer() - Start")

	for program_and_max_position := range channels_program_and_max_position {
		if has_stopped(&program_and_max_position.Program) {
			continue
		}
		tweets, max_position, err := tweets_fetch(
			settings, program_and_max_position.Program.Query, program_and_max_position.MaxPosition,
		)
		if err != nil {
			channels_program_and_max_position <- program_and_max_position
			continue
		}
		for _, tweet := range tweets {
			rest_api_tweet_insert(database, tweet)
		}
		if max_position != "" {
			var program_and_max_position = ProgramAndMaxPosition{
				Program:     program_and_max_position.Program,
				MaxPosition: max_position,
			}
			channels_program_and_max_position <- program_and_max_position
		}
	}

	fmt.Println("rest_api_consumer() - Stop")
}

func rest_api_producer(database *sqlx.DB, channels_program_and_max_position chan ProgramAndMaxPosition) {
	fmt.Println("rest_api_producer() - Start")

	tweets_delete(database)

	programs := programs_select(database)
	for _, program := range programs {
		var program_and_max_position = ProgramAndMaxPosition{
			Program:     program,
			MaxPosition: "",
		}
		channels_program_and_max_position <- program_and_max_position
	}

	fmt.Println("rest_api_producer() - Stop")
}

func rest_api_tweet_insert(database *sqlx.DB, tweet Tweet) {
	tweet_insert(database, &tweet)
}
