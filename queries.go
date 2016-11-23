package main

import (
	"github.com/jmoiron/sqlx"
)

func get_programs(database *sqlx.DB) []Program {
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
	rows, err := database.Queryx(query)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var program Program
		err = rows.StructScan(&program)
		if err != nil {
			panic(err)
		}
		programs = append(programs, program)
	}
	return programs
}

func get_screen_names(database *sqlx.DB) []string {
	var screen_names []string
	query := `
    SELECT twitter_user_screen_name
    FROM tweets
    WHERE
        twitter_user_tweets IS NOT NONE
        OR
        twitter_user_followers IS NOT NONE
        OR
        twitter_user_following IS NOT NONE
        OR
        twitter_user_created_at IS NOT NONE
    `
	err := database.Select(&screen_names, query)
	if err != nil {
		panic(err)
	}
	return screen_names
}

func get_tweeter(database *sqlx.DB, screen_name string) (*Tweeter, error) {
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
        twitter_user_tweets IS NOT NONE
        AND
        twitter_user_followers IS NOT NONE
        AND
        twitter_user_following IS NOT NONE
        AND
        twitter_user_created_at IS NOT NONE
    `
	row := database.QueryRowx(query, screen_name)
	err := row.StructScan(&tweeter)
	return &tweeter, err
}

func set_tweeter(database *sqlx.DB, tweeter *Tweeter) {
	database.NamedExec(
		`
        UPDATE tweets
        SET
            twitter_user_tweets = :tweets,
            twitter_user_followers = :followers,
            twitter_user_following = :following,
            twitter_user_created_at = :created_at
        WHERE twitter_user_screen_name = :screen_name
        `,
		tweeter,
	)
}

func set_tweet(database *sqlx.DB, tweet *Tweet) {
	database.NamedExec(
		`
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
        `,
		tweet,
	)
}
