package main

import (
	"github.com/jmoiron/sqlx"
	"strings"
)

func programs_select(database *sqlx.DB) []Program {
	var programs []Program
	query := `
    SELECT programs.queries_twitter, channels_programs.beginning_at, channels_programs.ending_at
    FROM programs
    INNER JOIN channels_programs ON channels_programs.program_id = programs.id
    WHERE
        channels_programs.beginning_at <= TIMEZONE('UTC', NOW())
        AND
        channels_programs.ending_at >= TIMEZONE('UTC', NOW())
    `
	err := database.Select(&programs, query)
	if err != nil {
		panic(err)
	}
	return programs
}

func screen_names_select(database *sqlx.DB) []string {
	var screen_names []string
	query := `
    SELECT DISTINCT twitter_user_screen_name
    FROM tweets
    WHERE
        twitter_user_tweets IS NOT NULL
        OR
        twitter_user_followers IS NOT NULL
        OR
        twitter_user_following IS NOT NULL
        OR
        twitter_user_timestamp IS NOT NULL
    `
	err := database.Select(&screen_names, query)
	if err != nil {
		panic(err)
	}
	return screen_names
}

func tweet_insert(database *sqlx.DB, tweet *Tweet) {
	tweet.Text = strings.Replace(tweet.Text, "#", "HASHTAG", -1)
	query := `
    INSERT INTO tweets
        (
            twitter_id,
            twitter_text,
            twitter_retweets,
            twitter_timestamp,
            twitter_user_id,
            twitter_user_screen_name,
            twitter_user_name,
            twitter_user_profile_image_url,
            twitter_user_tweets,
            twitter_user_followers,
            twitter_user_following,
            twitter_user_timestamp
        )
        VALUES
        (
            :twitter_id,
            :twitter_text,
            :twitter_retweets,
            :twitter_timestamp,
            :twitter_user_id,
            :twitter_user_screen_name,
            :twitter_user_name,
            :twitter_user_profile_image_url,
            :twitter_user_tweets,
            :twitter_user_followers,
            :twitter_user_following,
            :twitter_user_timestamp
        )
    `
	database.NamedExec(query, tweet)
}

func tweeter_select(database *sqlx.DB, screen_name string) (*Tweeter, error) {
	var tweeter Tweeter
	query := `
    SELECT
        twitter_user_screen_name,
        twitter_user_tweets,
        twitter_user_followers,
        twitter_user_following,
        twitter_user_timestamp
    FROM tweets
    WHERE
        twitter_user_screen_name = ?
        AND
        twitter_user_tweets IS NOT NULL
        AND
        twitter_user_followers IS NOT NULL
        AND
        twitter_user_following IS NOT NULL
        AND
        twitter_user_timestamp IS NOT NULL
    `
	err := database.Get(&tweeter, query, screen_name)
	return &tweeter, err
}

func tweeter_update(database *sqlx.DB, tweeter *Tweeter) {
	query := `
    UPDATE tweets
    SET
        twitter_user_tweets = :tweets,
        twitter_user_followers = :followers,
        twitter_user_following = :following,
        twitter_user_timestamp = :timestamp
    WHERE twitter_user_screen_name = :screen_name
    `
	database.NamedExec(query, tweeter)
}
