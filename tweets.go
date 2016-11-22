package main

import (
	"fmt"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var settings *Settings
	settings = get_settings()

	oauth1_config := oauth1.NewConfig(settings.Twitter.ConsumerKey, settings.Twitter.ConsumerSecret)
	oauth1_token := oauth1.NewToken(settings.Twitter.AccessKey, settings.Twitter.AccessSecret)
	oauth1_client := oauth1_config.Client(oauth1.NoContext, oauth1_token)

	parameters := &twitter.StreamFilterParams{
		Track: []string{
			"#TrumpAltRightFilms",
			"#bfc530",
			"#1MillionGifs",
			"#TuesdayMotivation",
			"President John F. Kennedy",
			"#DemotivateASong",
			"Toa Baja",
			"Donald Trump Rages",
			"Johnthony Walker",
			"Blowing Rock",
		},
	}

	client := twitter.NewClient(oauth1_client)
	stream, err := client.Streams.Filter(parameters)
	if err != nil {
		panic(err)
	}

	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		fmt.Println(tweet.IDStr)
	}
	demux.HandleChan(stream.Messages)

	go demux.HandleChan(stream.Messages)

	channel := make(chan os.Signal)
	signal.Notify(channel, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println(<-channel)

	stream.Stop()
}
