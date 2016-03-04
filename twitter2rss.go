package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/feeds"
	"github.com/volker-fr/twitter2rss/config"
	"github.com/volker-fr/twitter2rss/feed"
	"github.com/volker-fr/twitter2rss/parser"
	"github.com/volker-fr/twitter2rss/version"

	"github.com/davecgh/go-spew/spew"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

var conf config.Config = config.LoadConfig()

func processAPIError(message string, error error) {
	fmt.Printf(message)
	spew.Dump(error)
}

func getRss() string {
	twitterConfig := oauth1.NewConfig(conf.ConsumerKey, conf.ConsumerSecret)
	token := oauth1.NewToken(conf.AccessToken, conf.AccessSecret)
	// OAuth1 http.Client will automatically authorize Requests
	httpClient := twitterConfig.Client(oauth1.NoContext, token)

	// Twitter Client
	client := twitter.NewClient(httpClient)

	// debugging & testing
	if conf.Debug {
		var tweetId int64 = 7654321
		tweet, _, err := client.Statuses.Show(tweetId, &twitter.StatusShowParams{})
		if err != nil {
			processAPIError("Couldn't load client.Statuses.Show: ", err)
			return ""
		}
		fmt.Println(parser.GetTweetUrl(*tweet))
		spew.Dump(tweet)
		fmt.Println(parser.ParseTweetText(*tweet))
		return ""
	}

	// Get timeline
	homeTimelineParams := &twitter.HomeTimelineParams{Count: conf.MaxTweets}
	tweets, _, err := client.Timelines.HomeTimeline(homeTimelineParams)
	if err != nil {
		processAPIError("Couldn't load HomeTimeline: ", err)
		return ""
	}

	// get either single tweets or combine multiple tweets of the same author together
	var rssFeed *feeds.Feed
	if conf.CombinedFeed == false {
		rssFeed = feed.CreateIndividualFeed(conf, tweets)
	} else {
		rssFeed = feed.CreateCombinedUserFeed(conf, tweets)
	}

	// Create feed
	atom, err := rssFeed.ToAtom()
	if err != nil {
		log.Fatal(err)
	}

	return atom
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", getRss())
}

func main() {
	fmt.Println("twitter2rss version", version.GetVersion())
	if conf.Debug {
		_ = getRss()
	}

	// TODO: add logging
	// TODO: add error handling in case the port is already in use
	http.HandleFunc("/", handler)
	http.ListenAndServe("127.0.0.1:8080", nil)
}
