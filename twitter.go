package main

import (
	"encoding/json"
	"gopkg.in/xmlpath.v2"
	"net/http"
	"strconv"
	"strings"
)

func tweets_fetch(settings *Settings, q string, max_position string) ([]Tweet, string, error) {
	var err error

	client := get_http_client(settings, true)

	request, err := http.NewRequest("GET", "https://twitter.com/i/search/timeline", nil)
	if err != nil {
		return []Tweet{}, "", err
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
		return []Tweet{}, "", err
	}

	defer response.Body.Close()

	var items_and_max_position ItemsAndMaxPosition
	json.NewDecoder(response.Body).Decode(&items_and_max_position)

	reader := strings.NewReader(items_and_max_position.Items)
	root, err := xmlpath.ParseHTML(reader)
	if err != nil {
		return []Tweet{}, "", err
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

		path = `.//small[contains(@class, "time")]/a/span/@data-time`
		xpath = xmlpath.MustCompile(path)
		value, ok = xpath.String(iter.Node())
		if ok {
			created_at := get_timestamp_from_integer(value)
			tweet.CreatedAt = created_at
		}

		path = `.//@data-item-id`
		xpath = xmlpath.MustCompile(path)
		value, ok = xpath.String(iter.Node())
		if ok {
			tweet.Id = value
		}

		tweet.Source = ""

		path = `.//p[contains(@class, "tweet-text")]/text()`
		xpath = xmlpath.MustCompile(path)
		value, ok = xpath.String(iter.Node())
		if ok {
			tweet.Text = value
		}

		path = `.//span[contains(@class, "ProfileTweet-action--retweet")]/span/@data-tweet-stat-count`
		xpath = xmlpath.MustCompile(path)
		value, ok = xpath.String(iter.Node())
		if ok {
			integer, err = strconv.Atoi(value)
			if err != nil {
				panic(err)
			}
			tweet.Retweets = integer
		}

		path = `.//div[contains(@class, "original-tweet")]/@data-user-id`
		xpath = xmlpath.MustCompile(path)
		value, ok = xpath.String(iter.Node())
		if ok {
			tweet.UserId = value
		}

		path = `.//div[contains(@class, "original-tweet")]/@data-name`
		xpath = xmlpath.MustCompile(path)
		value, ok = xpath.String(iter.Node())
		if ok {
			tweet.UserName = value
		}

		path = `.//img[contains(@class, "avatar js-action-profile-avatar")]/@src`
		xpath = xmlpath.MustCompile(path)
		value, ok = xpath.String(iter.Node())
		if ok {
			tweet.UserProfileImageURL = value
		}

		path = `.//div[contains(@class, "original-tweet")]/@data-screen-name`
		xpath = xmlpath.MustCompile(path)
		value, ok = xpath.String(iter.Node())
		if ok {
			tweet.UserScreenName = value
		}

		tweets = append(tweets, tweet)
	}

	return tweets, items_and_max_position.MaxPosition, nil
}
