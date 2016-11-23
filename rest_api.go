package main

import (
	"fmt"
)

func rest_api(settings *Settings) {
	//     goroutine workers 32(incoming)
	//         for queries in range channel
	//             if program has stopped
	//                 do nothing
	//             else
	//                 A:
	//                 tweets = get first page
	//                 insert
	//                     goroutine update tweeter also
	//                 if next page
	//                     go to else A:
	//                 else:
	//                     if program has stopped
	//                         do nothing
	//                     else
	//                         channel <- queries
	//     goroutine refresh(every 1 minute)
	//         get_queries
	//         channel <- queries

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

func rest_api_consumer() {
	// for queries in range channel
	//    if program has stopped
	//        do nothing
	//    else
	//        A:
	//        tweets = get first page
	//        insert
	//            goroutine update tweeter also
	//        if next page
	//            go to else A:
	//        else:
	//            if program has stopped
	//                do nothing
	//            else
	//                channel <- queries
}

func rest_api_producer() {
	// get_queries
	// channel <- queries
}
