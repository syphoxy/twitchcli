package main

import (
	"fmt"
	"github.com/mrshankly/go-twitch/twitch"
	"log"
	"net/http"
)

func main() {
	client := twitch.NewClient(&http.Client{})
	opt := &twitch.ListOptions{
		Game: "Dota 2",
	}

	games, err := client.Streams.List(opt)
	if err != nil {
		log.Fatal(err)
	}

	for _, s := range games.Streams {
		fmt.Printf("%6d %s: %s\n", s.Viewers, s.Channel.DisplayName, s.Channel.Status)
	}
}
