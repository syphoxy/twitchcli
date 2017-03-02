package main

import (
	"flag"
	"fmt"
	"github.com/syphoxy/go-twitch/twitch"
	"log"
	"net/http"
	"os"
)

func main() {
	client := twitch.NewClient(&http.Client{})

	var game string
	var limit int

	flag.StringVar(&game, "g", "", "Game")
	flag.IntVar(&limit, "l", 25, "Limit")
	flag.Parse()

	if limit <= 0 {
		os.Exit(0)
	}

	if limit > 100 {
		fmt.Println("Limit must be a value between 0 and 100")
		os.Exit(1)
	}

	options := twitch.ListOptions{}
	options.Limit = limit

	if game != "" {
		options.Game = game
	}

	games, err := client.Streams.List(&options)

	if err != nil {
		log.Fatal(err)
	}

	for _, s := range games.Streams {
		fmt.Printf("%8d %25s %s\n", s.Viewers, s.Channel.DisplayName, s.Channel.Status)
	}
}
