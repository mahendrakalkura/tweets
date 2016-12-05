package main

import (
	"github.com/jmoiron/sqlx"
	"strings"
)

func programs_select(database *sqlx.DB) []Program {
	var programs []Program
	query := `
    SELECT programs.query, channels_programs.beginning_at, channels_programs.ending_at
    FROM programs
    INNER JOIN channels_programs ON channels_programs.program_id = programs.id
    WHERE
        programs.status = TRUE
        AND
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

func tweet_insert(database *sqlx.DB, tweet *Tweet) {
	tweet.Text = strings.Replace(tweet.Text, "#", "HASHTAG", -1)
	query := `
    INSERT INTO tweets
        (
            twitter_created_at,
            twitter_id,
            twitter_source,
            twitter_text,
            twitter_retweets,
            twitter_user_id,
            twitter_user_name,
            twitter_user_profile_image_url,
            twitter_user_screen_name
        )
        VALUES
        (
            :twitter_created_at,
            :twitter_id,
            :twitter_source,
            :twitter_text,
            :twitter_retweets,
            :twitter_user_id,
            :twitter_user_name,
            :twitter_user_profile_image_url,
            :twitter_user_screen_name
        )
    `
	database.NamedExec(query, tweet)
}
