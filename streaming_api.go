package main

import (
	"fmt"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/jmoiron/sqlx"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"
)

func streaming_api(settings *Settings, database *sqlx.DB) {
	fmt.Println("streaming_api() - Start")

	channels_track := make(chan []string)
	channels_signal := make(chan os.Signal)
	channels_exit := make(chan bool)

	go streaming_api_consumer(settings, database, channels_track, channels_exit)

	go streaming_api_producer(database, channels_track)

	ticker := time.NewTicker(60 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				go streaming_api_producer(database, channels_track)
			}
		}
	}()

	signal.Notify(channels_signal, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)

	<-channels_signal

	close(channels_track)

	channels_exit <- true

	<-channels_exit

	fmt.Println("streaming_api() - Stop")
}

func streaming_api_consumer(
	settings *Settings, database *sqlx.DB, channels_track chan []string, channels_exit chan bool,
) {
	fmt.Println("streaming_api_consumer() - Start")

	tracks_old := []string{}
	tracks_new := []string{}

	channels_stop := make(chan bool)

	for track := range channels_track {
		tracks_old = tracks_new
		tracks_new = track

		if reflect.DeepEqual(tracks_old, tracks_new) {
			continue
		}

		if len(tracks_old) != 0 {
			channels_stop <- true
			<-channels_stop
		}

		if len(tracks_new) == 0 {
			continue
		}

		go streaming_api_consumer_stream(settings, database, tracks_new, channels_stop)
	}

	<-channels_exit

	if len(tracks_new) != 0 {
		channels_stop <- true
		<-channels_stop
	}

	fmt.Println("streaming_api_consumer() - Stop")

	channels_exit <- true
}

func streaming_api_producer(database *sqlx.DB, channels_track chan []string) {
	programs := get_programs(database)
	track := get_track(programs)
	channels_track <- track
}

func streaming_api_consumer_stream(settings *Settings, database *sqlx.DB, track []string, channels_stop chan bool) {
	fmt.Println("streaming_api_consumer_stream() - Start")

	client := get_twitter_client(
		settings.Twitter.ConsumerKey,
		settings.Twitter.ConsumerSecret,
		settings.Twitter.AccessKey,
		settings.Twitter.AccessSecret,
	)

	parameters := get_twitter_parameters(track)

	stream, err := client.Streams.Filter(parameters)
	if err != nil {
		panic(err)
	}

	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(twitter_tweet *twitter.Tweet) {
		go streaming_api_set_tweet(database, twitter_tweet)
	}

	go demux.HandleChan(stream.Messages)

	<-channels_stop

	stream.Stop()

	fmt.Println("streaming_api_consumer_stream() - Stop")

	channels_stop <- true
}

func streaming_api_set_tweet(database *sqlx.DB, twitter_tweet *twitter.Tweet) {
	created_at := get_timestamp_from_string(twitter_tweet.CreatedAt)
	user_created_at := get_timestamp_from_string(twitter_tweet.User.CreatedAt)
	tweet := &Tweet{
		Id:                  twitter_tweet.IDStr,
		Text:                twitter_tweet.Text,
		Retweets:            twitter_tweet.RetweetCount,
		CreatedAt:           created_at,
		UserId:              twitter_tweet.User.IDStr,
		UserScreenName:      twitter_tweet.User.ScreenName,
		UserName:            twitter_tweet.User.Name,
		UserProfileImageURL: twitter_tweet.User.ProfileImageURL,
		UserTweets:          &twitter_tweet.User.StatusesCount,
		UserFollowers:       &twitter_tweet.User.FollowersCount,
		UserFollowing:       &twitter_tweet.User.FriendsCount,
		UserCreatedAt:       &user_created_at,
	}
	set_tweet(database, tweet)
}
