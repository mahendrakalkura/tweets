package main

import (
	"flag"
)

func main() {
	action := flag.String("action", "", "")

	flag.Parse()

	var settings *Settings

	settings = get_settings()

	if *action == "streaming-api" {
		streaming_api(settings)
	}

	if *action == "rest-api" {
		rest_api(settings)
	}
}
