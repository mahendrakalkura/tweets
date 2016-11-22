package main

import (
    "encoding/json"
    "fmt"
    "gopkg.in/xmlpath.v2"
    "net/http"
    "strconv"
    "strings"
)

func main() {
    var settings *Settings
    settings = get_settings()

    var tweets_all []Tweet
    var tweets_some []Tweet
    var max_position string

    max_position = ""
    for {
        tweets_some, max_position = get_tweets(settings, "donald trump", max_position)
        if max_position == "" {
            break
        }
        tweets_all = append(tweets_all, tweets_some...)
    }

    fmt.Println(tweets_all)
}

func get_tweets(settings *Settings, q string, max_position string) ([]Tweet, string) {
    var err error

    client := get_client(settings, true)

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

    var body struct {
        Items       string `json:"items_html"`
        MaxPosition string `json:"min_position"`
    }
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
