package main

import (
    "flag"
    "fmt"
    "log"
    "os"
    "time"
    "net/http"
    "sort"

    "github.com/coreos/pkg/flagutil"
    "github.com/dghubble/go-twitter/twitter"
    "github.com/dghubble/oauth1"
    "github.com/davecgh/go-spew/spew"
    "github.com/gorilla/feeds"
)

type replaceObject struct {
    from int
    to int
    replacement string
}

type ReplacementList []replaceObject

func (dl ReplacementList) Len() int {
   return len(dl)
}

func (dl ReplacementList) Swap(i, j int) {
    dl[i], dl[j] = dl[j], dl[i]
}

func (dl ReplacementList) Less(i, j int) bool {
    return dl[i].from < dl[j].from
}



// Parse the text, identify twitter shortened URLs and replace them
func parseTweetText(tweet twitter.Tweet) string {
    text := tweet.Text
    var replacements ReplacementList

    // Go through each URL object and replace it with a link and correct text
    for _, url := range tweet.Entities.Urls {
        replacement := "<a href='" + url.ExpandedURL + "'>" + url.DisplayURL + "</a>"
        from := url.Indices[0]
        to := url.Indices[1]
        replacements = append(replacements, replaceObject{ from, to, replacement })
    }

    // Go through each Media object and replace it the link in the text with it
    for _, media := range tweet.Entities.Media {
        var mediaUrl, replacement string
        if media.Type != "photo" {
            fmt.Println("media.Type not photo")
            spew.Dump(tweet)
            replacement = "unsupported_mediatype"
        } else if media.Type == "photo" {
            // or maybe we should use media.URLEntity.ExpandedURL?
            if len(media.MediaURLHttps) != 0 {
               mediaUrl = media.MediaURLHttps
            } else {
               mediaUrl = media.MediaURL
            }
            replacement = "<br><img src='" + mediaUrl + "'><br>"
        }
        from := media.Indices[0]
        to := media.Indices[1]
        replacements = append(replacements, replaceObject{ from, to, replacement })
    }

    sort.Sort(replacements)

    // replacement is sorted, start from the end, since we change the length of the string
    //fmt.Println(text)
    for i := len(replacements) - 1; i >= 0; i-- {
      text = text[:replacements[i].from] + replacements[i].replacement + text[replacements[i].to:]
      //fmt.Println(text)
    }
    //fmt.Println("---\n")

    // Does this tweet quote another tweet?
    if tweet.QuotedStatus != nil {
        quotedTweet := *tweet.QuotedStatus
        url := "https://twitter.com/" + quotedTweet.User.ScreenName + "/status/" + quotedTweet.IDStr
        header := "<a href=\"" + url + "\">" + quotedTweet.User.Name + "</a><br>\n"
        quotedText := parseTweetText(quotedTweet)
        text +=  "\n<blockquote>\n" + header + quotedText + "\n</blockquote>\n"
    }

    return text
}

func getRss() string {
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

    /* // debugging & testing
    tweetId := 12345
    tweet, _, _ := client.Statuses.Show(tweetId, twitter.StatusShowParams{})
    fmt.Println("https://twitter.com/" + tweet.User.ScreenName + "/status/" + tweet.IDStr)
    println(parseTweetText(*tweet))
    return "" */

    for _,tweet := range tweets {
        if *debug { spew.Dump(tweet) }

        t, _ := time.Parse(time.RubyDate, tweet.CreatedAt)

        url := "https://twitter.com/" + tweet.User.ScreenName + "/status/" + tweet.IDStr
        if *debug { fmt.Println(url) }

        item := &feeds.Item{
          Title:       fmt.Sprintf("%s: %s...", tweet.User.Name, tweet.Text[:10] ),
          Link:        &feeds.Link{Href: url},
          Description: tweet.User.Name + ": " + parseTweetText(tweet),
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
    fmt.Fprintf(w, "%s", getRss())
}

func main() {
    //_ = getRss()

    http.HandleFunc("/", handler)
    http.ListenAndServe("127.0.0.1:8080", nil)
}
