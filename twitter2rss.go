package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/coreos/pkg/flagutil"
	"github.com/davecgh/go-spew/spew"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/gorilla/feeds"
)

type replaceObject struct {
	from        int
	to          int
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

func getTweetUrl(tweet twitter.Tweet) string {
	return "https://twitter.com/" + tweet.User.ScreenName + "/status/" + tweet.IDStr
}

func getTweetText(id int64, client *twitter.Client) string {
	tweet, _, _ := client.Statuses.OEmbed(id, &twitter.StatusOEmbedParams{})

	return tweet.HTML
}

// Parse the text, identify twitter shortened URLs and replace them
func parseTweetText(tweet twitter.Tweet) string {
	text := tweet.Text
	var replacements ReplacementList

	// Special case, if it's retweeted then the URL placement might not be correct
	// and the tweet can also contain cut off text.
	// The Twitter timeline also doesn't show the RT message but the retweeted Tweet.
	if tweet.RetweetedStatus != nil {
		text = parseTweetText(*tweet.RetweetedStatus)
		return "<a href=\"" + getTweetUrl(tweet) + "\">" + tweet.User.Name + "</a>: RT @" + tweet.RetweetedStatus.User.ScreenName + ":<br>\n" + text
	}

	// Go through each URL object and replace it with a link and correct text
	for _, url := range tweet.Entities.Urls {
		replacement := "<a href='" + url.ExpandedURL + "'>" + url.DisplayURL + "</a>"
		from := url.Indices[0]
		to := url.Indices[1]
		replacements = append(replacements, replaceObject{from, to, replacement})
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
		replacements = append(replacements, replaceObject{from, to, replacement})
	}

	sort.Sort(replacements)

	// replacement is sorted, start from the end, since we change the length of the string
	//fmt.Println(text)
	for i := len(replacements) - 1; i >= 0; i-- {
		// cut text after character and not byte by converting the string to a rune and back to a string
		text = string([]rune(text)[:replacements[i].from]) + replacements[i].replacement + string([]rune(text)[replacements[i].to:])
		//fmt.Println(text)
	}
	//fmt.Println("---\n")

	// Does this tweet quote another tweet?
	if tweet.QuotedStatus != nil {
		quotedTweet := *tweet.QuotedStatus
		header := "<a href=\"" + getTweetUrl(quotedTweet) + "\">" + quotedTweet.User.Name + "</a><br>\n"
		quotedText := parseTweetText(quotedTweet)
		text += "\n<blockquote>\n" + header + quotedText + "\n</blockquote>\n"
	}

	return "<a href=\"" + getTweetUrl(tweet) + "\">" + tweet.User.Name + "</a>:<br>\n" + text
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
		var tweetId int64 = 1234
		tweet, _, _ := client.Statuses.OEmbed(tweetId, &twitter.StatusOEmbedParams{})
		spew.Dump(tweet)
	    //tweet, _, _ := client.Statuses.Show(tweetId, &twitter.StatusShowParams{})
		//fmt.Println("https://twitter.com/" + tweet.User.ScreenName + "/status/" + tweet.IDStr)
	    //println(parseTweetText(*tweet))
		return "" */

	for _, tweet := range tweets {
		if *debug {
			spew.Dump(tweet)
		}

		t, _ := time.Parse(time.RubyDate, tweet.CreatedAt)

		item := &feeds.Item{
			Title: fmt.Sprintf("%s: %s...", tweet.User.Name, tweet.Text[:10]),
			Link:  &feeds.Link{Href: getTweetUrl(tweet)},
			//Description: parseTweetText(tweet),
			Description: getTweetText(tweet.ID, client),
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

	if *debug {
		fmt.Println(atom, "\n")
	}

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
