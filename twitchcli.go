package main

import (
	"flag"
	"fmt"
	"github.com/syphoxy/go-twitch/twitch"
	"log"
	"math"
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

	displayNameLen := 0
	for _, s := range games.Streams {
		length := len(s.Channel.DisplayName)
		if length > displayNameLen {
			displayNameLen = length
		}
	}

	viewersLen := int(math.Ceil(math.Log(float64(games.Streams[0].Viewers)) / math.Log(10)))
	streamFmt := fmt.Sprintf("  %%-%ds %%%dd %%s\n", displayNameLen, viewersLen)

	for _, s := range games.Streams {
		fmt.Printf(streamFmt, s.Channel.DisplayName, s.Viewers, s.Channel.Status)
	}
}
