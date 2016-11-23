package main

import (
	"github.com/jmoiron/sqlx"
)

func programs_select(database *sqlx.DB) ([]Program, error) {
	var programs []Program
	query := `
    SELECT programs.query, channels_programs.beginning_at, channels_programs.ending_at
    FROM programs
    INNER JOIN channels_programs ON channels_programs.program_id = programs.id
    WHERE
        channels_programs.beginning_at <= TIMEZONE('UTC', NOW())
        AND
        channels_programs.ending_at >= TIMEZONE('UTC', NOW())
    `
	err := database.Select(&programs, query)
	return programs, err
}

func screen_names_select(database *sqlx.DB) ([]string, error) {
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
        twitter_user_created_at IS NOT NULL
    `
	err := database.Select(&screen_names, query)
	return screen_names, err
}

func tweet_insert(database *sqlx.DB, tweet *Tweet) {
	query := `
    INSERT INTO tweets
        (
            twitter_id,
            twitter_text,
            twitter_retweets,
            twitter_created_at,
            twitter_user_id,
            twitter_user_screen_name,
            twitter_user_name,
            twitter_user_profile_image_url,
            twitter_user_tweets,
            twitter_user_followers,
            twitter_user_following,
            twitter_user_created_at
        )
        VALUES
        (
            :twitter_id,
            :twitter_text,
            :twitter_retweets,
            :twitter_created_at,
            :twitter_user_id,
            :twitter_user_screen_name,
            :twitter_user_name,
            :twitter_user_profile_image_url,
            :twitter_user_tweets,
            :twitter_user_followers,
            :twitter_user_following,
            :twitter_user_created_at
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
        twitter_user_created_at
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
        twitter_user_created_at IS NOT NULL
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
        twitter_user_created_at = :created_at
    WHERE twitter_user_screen_name = :screen_name
    `
	database.NamedExec(query, tweeter)
}
