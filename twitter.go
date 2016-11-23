package main

import (
	"encoding/json"
	"fmt"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"gopkg.in/xmlpath.v2"
	"net/http"
	"strconv"
	"strings"
)

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

func get_tweets(settings *Settings, q string, max_position string) ([]Tweet, string) {
	var err error

	client := get_http_client(settings, true)

	request, err := http.NewRequest("GET", "https://twitter.com/i/search/timeline", nil)
	if err != nil {
		panic(err)
	}

	request.Header.Add("accept", "application/json")
	request.Header.Add("x-push-state-request", "true")
	request.Header.Add("x-requested-with", "XMLHttpRequest")

	query := request.URL.Query()
	query.Add("composed_count", "0")
	query.Add("include_available_features", "1")
	query.Add("include_entities", "1")
	query.Add("include_new_items_bar", "true")
	query.Add("latent_count", "0")
	query.Add("max_position", max_position)
	query.Add("q", q)
	query.Add("src", "typd")
	query.Add("vertical", "news")
	request.URL.RawQuery = query.Encode()

	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}

	var body Body
	json.NewDecoder(response.Body).Decode(&body)

	reader := strings.NewReader(body.Items)
	root, err := xmlpath.ParseHTML(reader)
	if err != nil {
		panic(err)
	}

	var path string
	var xpath *xmlpath.Path
	var iter *xmlpath.Iter
	var value string
	var ok bool
	var integer int
	var tweets []Tweet
	var tweet Tweet

	path = `//li[@data-item-type="tweet"]`
	xpath = xmlpath.MustCompile(path)
	iter = xpath.Iter(root)
	for iter.Next() {
		tweet = Tweet{}

		path = `.//@data-item-id`
		xpath = xmlpath.MustCompile(path)
		value, ok = xpath.String(iter.Node())
		if ok {
			tweet.Id = value
		}

		path = ".//div[contains(@class, \"js-tweet-text-container\")]"
		path += "/p[contains(@class, \"tweet-text\")]"
		path += "/text()"
		xpath = xmlpath.MustCompile(path)
		value, ok = xpath.String(iter.Node())
		if ok {
			tweet.Text = value
		}

		path = ".//span[contains(@class, \"ProfileTweet-action--retweet\")]"
		path += "/span[contains(@class, \"ProfileTweet-actionCount\")]"
		path += "/@data-tweet-stat-count"
		xpath = xmlpath.MustCompile(path)
		value, ok = xpath.String(iter.Node())
		if ok {
			integer, err = strconv.Atoi(value)
			if err != nil {
				panic(err)
			}
			tweet.Retweets = integer
		}

		path = ".//div[contains(@class, \"stream-item-header\")]"
		path += "/small[contains(@class, \"time\")]"
		path += "/a"
		path += "/span"
		path += "/@data-time"
		xpath = xmlpath.MustCompile(path)
		value, ok = xpath.String(iter.Node())
		if ok {
			unix := get_unix(value)
			tweet.CreatedAt = unix
		}

		path = ".//div[contains(@class, \"original-tweet\")]"
		path += "/@data-user-id"
		xpath = xmlpath.MustCompile(path)
		value, ok = xpath.String(iter.Node())
		if ok {
			tweet.UserId = value
		}

		path = ".//div[contains(@class, \"original-tweet\")]"
		path += "/@data-screen-name"
		xpath = xmlpath.MustCompile(path)
		value, ok = xpath.String(iter.Node())
		if ok {
			tweet.UserScreenName = value
		}

		path = ".//div[contains(@class, \"original-tweet\")]"
		path += "/@data-name"
		xpath = xmlpath.MustCompile(path)
		value, ok = xpath.String(iter.Node())
		if ok {
			tweet.UserName = value
		}

		path = ".//img[contains(@class, \"avatar js-action-profile-avatar\")]"
		path += "/@src"
		xpath = xmlpath.MustCompile(path)
		value, ok = xpath.String(iter.Node())
		if ok {
			tweet.UserProfileImageUrl = value
		}

		// UserTweets

		// UserFollowers

		// UserFollowing

		// UserCreatedAt

		fmt.Println(tweet.Id)

		tweets = append(tweets, tweet)
	}

	return tweets, body.MaxPosition
}
