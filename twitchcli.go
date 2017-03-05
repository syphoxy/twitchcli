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
	options := twitch.ListOptions{
		Limit: limit,
	}
	if game != "" {
		options.Game = game
	}
	followed, err := client.Streams.Followed(&options)
	if err != nil {
		log.Fatal(err)
	}
	streams, err := client.Streams.List(&options)
	if err != nil {
		log.Fatal(err)
	}
	RenderList(followed.Streams)
	fmt.Println("")
	RenderList(streams.Streams)
}

func RenderList(list []twitch.StreamS) {
	nameLen := 0
	for _, s := range list {
		length := len(s.Channel.Name)
		if length > nameLen {
			nameLen = length
		}
	}
	viewersLen := int(math.Ceil(math.Log(float64(list[0].Viewers)) / math.Log(10)))
	streamFmt := fmt.Sprintf("%%-%ds %%%dd %%s\n", nameLen, viewersLen)
	for _, s := range list {
		fmt.Printf(streamFmt, s.Channel.Name, s.Viewers, s.Channel.Status)
	}
}
