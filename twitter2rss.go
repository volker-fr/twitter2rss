package main

import (
    "flag"
    "fmt"
    "log"
    "os"
    "time"

    "github.com/coreos/pkg/flagutil"
    "github.com/dghubble/go-twitter/twitter"
    "github.com/dghubble/oauth1"
	  "github.com/davecgh/go-spew/spew"
    "github.com/gorilla/feeds"
)

func main() {
    flags := flag.NewFlagSet("user-auth", flag.ExitOnError)
    consumerKey := flags.String("consumer-key", "", "Twitter Consumer Key")
    consumerSecret := flags.String("consumer-secret", "", "Twitter Consumer Secret")
    accessToken := flags.String("access-token", "", "Twitter Access Token")
    accessSecret := flags.String("access-secret", "", "Twitter Access Secret")
    debug := flags.Bool("debug", false, "Debug")

    flags.Parse(os.Args[1:])
    flagutil.SetFlagsFromEnv(flags, "TWITTER")

    if *consumerKey == "" || *consumerSecret == "" || *accessToken == "" || *accessSecret == "" {
        log.Fatal("Consumer key/secret and Access token/secret required")
    }

    config := oauth1.NewConfig(*consumerKey, *consumerSecret)
    token := oauth1.NewToken(*accessToken, *accessSecret)
    // OAuth1 http.Client will automatically authorize Requests
    httpClient := config.Client(oauth1.NoContext, token)

    // Twitter Client
    client := twitter.NewClient(httpClient)

    now := time.Now()

    // rss feed
    feed := &feeds.Feed{
      Title:       "Twitter Home Timeline",
      Description: "Twitter Home Timeline RSS Feed",
      Author:      &feeds.Author{"My name", "my@email"},
      Created:     now,
    }
    feed.Items = []*feeds.Item{}

    // Get timeline
    homeTimelineParams := &twitter.HomeTimelineParams{Count: 5}
    tweets, _, _ := client.Timelines.HomeTimeline(homeTimelineParams)

    for _,tweet := range tweets {
        if *debug {
          spew.Dump(tweet)
        }

        item := &feeds.Item{
          Title:       "Twitter content 1",
          Link:        &feeds.Link{Href: "http://tobedefined"},
          Description: tweet.Text,
          Author:      &feeds.Author{tweet.User.Name, tweet.User.ScreenName},
          Created:     now, //tweet.CreatedAt,
        }
        feed.Items = append(feed.Items, item)

    		fmt.Println(tweet.IDStr)
    		fmt.Println(tweet.CreatedAt)
    		fmt.Println(tweet.Text)
        fmt.Printf("%v (%v)\n", tweet.User.Name, tweet.User.ScreenName)
        fmt.Println()
    }

    // Create feed
    atom, err := feed.ToAtom()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(atom, "\n")
//spew.Dump(feed)

}
