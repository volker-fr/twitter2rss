package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"sort"
	"time"

	"github.com/volker-fr/twitter2rss/config"

	"github.com/davecgh/go-spew/spew"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/gorilla/feeds"
)

var conf config.Config = config.LoadConfig()

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

func processAPIError(message string, error error) {
	fmt.Printf(message)
	spew.Dump(error)
}

func convertTwitterTime(timestring string) time.Time {
	t, _ := time.Parse(time.RubyDate, timestring)
	return t
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
	if tweet.ExtendedEntities != nil && tweet.ExtendedEntities.Media != nil {
		var mediaFrom, mediaTo int
		// Media objects always have the same indices for replacement for multiple objects
		// therefore we append the data to the existing string
		var mediaReplacement string
		for _, media := range tweet.ExtendedEntities.Media {
			var mediaUrl string
			if media.Type != "photo" {
				fmt.Println("media.Type not photo")
				spew.Dump(tweet)
				mediaReplacement += "unsupported_mediatype"
			} else if media.Type == "photo" {
				// or maybe we should use media.URLEntity.ExpandedURL?
				if len(media.MediaURLHttps) != 0 {
					mediaUrl = media.MediaURLHttps
				} else {
					mediaUrl = media.MediaURL
				}
				mediaReplacement += "<br><img src='" + mediaUrl + "'><br>"
			}
			mediaFrom = media.Indices[0]
			mediaTo = media.Indices[1]
			replacements = append(replacements, replaceObject{mediaFrom, mediaTo, mediaReplacement})
		}
	}

	sort.Sort(replacements)

	// replacement is sorted, start from the end, since we change the length of the string
	for i := len(replacements) - 1; i >= 0; i-- {
		// cut text after character and not byte by converting the string to a rune and back to a string
		text = string([]rune(text)[:replacements[i].from]) + replacements[i].replacement + string([]rune(text)[replacements[i].to:])
	}

	// Does this tweet quote another tweet?
	if tweet.QuotedStatus != nil {
		quotedTweet := *tweet.QuotedStatus
		header := "<a href=\"" + getTweetUrl(quotedTweet) + "\">" + quotedTweet.User.Name + "</a><br>\n"
		quotedText := parseTweetText(quotedTweet)
		text += "\n<blockquote>\n" + header + quotedText + "\n</blockquote>\n"
	}

	footer := "<p><a href=\"" + getTweetUrl(tweet) + "\">" + tweet.User.Name + "</a> @ " + convertTwitterTime(tweet.CreatedAt).Format(time.RFC1123) + "\n"

	return text + footer
}

func isTweetFiltered(tweet twitter.Tweet, conf config.Config) bool {
	if len(conf.IgnoreSource) > 0 {
		for _, searchString := range conf.IgnoreSource {
			re := regexp.MustCompile(searchString)
			if re.MatchString(tweet.Source) {
				fmt.Printf("Source filter regex matched for tweet %s: %s\n", tweet.IDStr, searchString)
				return true
			}
		}
	}

	if len(conf.IgnoreText) > 0 {
		// parse it so we can also filter the URLs in a tweet
		// we still parse this tweet twice... once here, once at the feed creation :(
		parsedTweetText := parseTweetText(tweet)
		for _, searchString := range conf.IgnoreText {
			re := regexp.MustCompile(searchString)
			if re.MatchString(parsedTweetText) {
				fmt.Printf("Text filter regex matched for tweet %s: %s\n", tweet.IDStr, searchString)
				return true
			}
		}
	}

	return false
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
		Author:      &feeds.Author{"Twitter2RSS", "lists.volker@gmail.com"},
		Link:        &feeds.Link{Href: "http://github.com:volker-fr/twitter2rss/"},
		Created:     now,
	}
	feed.Items = []*feeds.Item{}

	/* // debugging & testing
		//var tweetId int64 = 1234
	    tweet, _, err := client.Statuses.Show(tweetId, &twitter.StatusShowParams{})
		if err != nil {
			processAPIError("Couldn't load client.Statuses.Show: ", err)
			return ""
		}
		fmt.Println("https://twitter.com/" + tweet.User.ScreenName + "/status/" + tweet.IDStr)
		spew.Dump(tweet)
		println(parseTweetText(*tweet))
		return "" */

	// Get timeline
	homeTimelineParams := &twitter.HomeTimelineParams{Count: 50}
	tweets, _, err := client.Timelines.HomeTimeline(homeTimelineParams)
	if err != nil {
		processAPIError("Couldn't load HomeTimeline: ", err)
		return ""
	}

	for _, tweet := range tweets {
		if isTweetFiltered(tweet, conf) {
			continue
		}

		titleLimit := 40
		if len(tweet.Text) < 40 {
			titleLimit = len(tweet.Text)
		}
		item := &feeds.Item{
			// TODO: check if slicing a string with non ascii chars will fail/scramble the text
			Title:       fmt.Sprintf("%s: %s...", tweet.User.Name, tweet.Text[:titleLimit]),
			Link:        &feeds.Link{Href: getTweetUrl(tweet)},
			Description: parseTweetText(tweet),
			Author:      &feeds.Author{tweet.User.Name, tweet.User.ScreenName},
			Created:     convertTwitterTime(tweet.CreatedAt),
			Id:          tweet.IDStr,
		}
		feed.Add(item)
	}

	// Create feed
	atom, err := feed.ToAtom()
	if err != nil {
		log.Fatal(err)
	}

	if conf.Debug {
		fmt.Println(atom, "\n")
	}

	return atom
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", getRss())
}

func main() {
	//_ = getRss()

	// TODO: add logging
	// TODO: add error handling in case the port is already in use
	http.HandleFunc("/", handler)
	http.ListenAndServe("127.0.0.1:8080", nil)
}
