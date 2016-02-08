package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/volker-fr/twitter2rss/config"
	"github.com/volker-fr/twitter2rss/filter"
	"github.com/volker-fr/twitter2rss/parser"

	"github.com/davecgh/go-spew/spew"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/gorilla/feeds"
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

	now := time.Now()

	// rss feed
	feed := &feeds.Feed{
		Title:       "Twitter Home Timeline",
		Description: "Twitter Home Timeline RSS Feed",
		Author:      &feeds.Author{Name: "Twitter2RSS", Email: "lists.volker@gmail.com"},
		Link:        &feeds.Link{Href: "http://github.com:volker-fr/twitter2rss/"},
		Created:     now,
	}
	feed.Items = []*feeds.Item{}

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
		println(parser.ParseTweetText(*tweet))
		return ""
	}

	// Get timeline
	homeTimelineParams := &twitter.HomeTimelineParams{Count: 50}
	tweets, _, err := client.Timelines.HomeTimeline(homeTimelineParams)
	if err != nil {
		processAPIError("Couldn't load HomeTimeline: ", err)
		return ""
	}

	for _, tweet := range tweets {
		parsedTweetText := parser.ParseTweetText(tweet)

		if filter.IsTweetFiltered(tweet, conf, parsedTweetText) {
			continue
		}

		titleLimit := 40
		if len(tweet.Text) < 40 {
			titleLimit = len(tweet.Text)
		}
		item := &feeds.Item{
			// TODO: check if slicing a string with non ascii chars will fail/scramble the text
			Title:       fmt.Sprintf("%s: %s...", tweet.User.Name, tweet.Text[:titleLimit]),
			Link:        &feeds.Link{Href: parser.GetTweetUrl(tweet)},
			Description: parsedTweetText,
			Author:      &feeds.Author{Name: tweet.User.Name, Email: tweet.User.ScreenName},
			Created:     parser.ConvertTwitterTime(tweet.CreatedAt),
			Id:          tweet.IDStr,
		}
		feed.Add(item)
	}

	// Create feed
	atom, err := feed.ToAtom()
	if err != nil {
		log.Fatal(err)
	}

	return atom
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", getRss())
}

func main() {
	if conf.Debug {
		_ = getRss()
	}

	// TODO: add logging
	// TODO: add error handling in case the port is already in use
	http.HandleFunc("/", handler)
	http.ListenAndServe("127.0.0.1:8080", nil)
}
