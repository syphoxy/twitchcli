package main

import (
	"fmt"
	"github.com/syphoxy/go-twitch/twitch"
	"log"
	"net/http"
)

func main() {
	client := twitch.NewClient(&http.Client{})

	games, err := client.Streams.List(&twitch.ListOptions{
		Game: "Dota 2",
	})

	if err != nil {
		log.Fatal(err)
	}

	for _, s := range games.Streams {
		fmt.Printf("%8d %25s %s\n", s.Viewers, s.Channel.DisplayName, s.Channel.Status)
	}
}
