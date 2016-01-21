package main

import (
    "flag"
    "fmt"
    "log"
    "os"
    "time"
    "net/http"

    "github.com/coreos/pkg/flagutil"
    "github.com/dghubble/go-twitter/twitter"
    "github.com/dghubble/oauth1"
    "github.com/davecgh/go-spew/spew"
    "github.com/gorilla/feeds"
)


func get_rss() string {
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
      Author:      &feeds.Author{"Twitter2RSS", "lists.volker@gmail.com"},
      Link:        &feeds.Link{Href: "http://github.com:volker-fr/twitter2rss/"},
      Created:     now,
    }
    feed.Items = []*feeds.Item{}

    // Get timeline
    homeTimelineParams := &twitter.HomeTimelineParams{Count: 50}
    tweets, _, _ := client.Timelines.HomeTimeline(homeTimelineParams)

    for _,tweet := range tweets {
        if *debug { spew.Dump(tweet) }

        t, _ := time.Parse(time.RubyDate, tweet.CreatedAt)

        url := "https://twitter.com/" + tweet.User.ScreenName + "/status/" + tweet.IDStr

        item := &feeds.Item{
          Title:       fmt.Sprintf("%s: %s...", tweet.User.Name, tweet.Text[:10] ),
          Link:        &feeds.Link{Href: url},
          Description: tweet.Text,
          Author:      &feeds.Author{tweet.User.Name, tweet.User.ScreenName},
          Created:     t,
          Id:          tweet.IDStr,
        }
        feed.Add(item)
    }

    // Create feed
    atom, err := feed.ToAtom()
    if err != nil {
        log.Fatal(err)
    }

    if *debug { fmt.Println(atom, "\n") }

    return atom
}

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "%s", get_rss())
}

func main() {
    http.HandleFunc("/", handler)
    http.ListenAndServe(":8080", nil)
}
