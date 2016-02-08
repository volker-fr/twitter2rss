package parser

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/dghubble/go-twitter/twitter"
)

func GetTweetUrl(tweet twitter.Tweet) string {
	return "https://twitter.com/" + tweet.User.ScreenName + "/status/" + tweet.IDStr
}

func ConvertTwitterTime(timestring string) time.Time {
	t, _ := time.Parse(time.RubyDate, timestring)
	return t
}

// Parse the text, identify twitter shortened URLs and replace them
func ParseTweetText(tweet twitter.Tweet) string {
	text := tweet.Text
	var replacements ReplacementList

	// Special case, if it's retweeted then the URL placement might not be correct
	// and the tweet can also contain cut off text.
	// The Twitter timeline also doesn't show the RT message but the retweeted Tweet.
	if tweet.RetweetedStatus != nil {
		text = ParseTweetText(*tweet.RetweetedStatus)
		return text + "<br>\nvia RT from <a href=\"" + GetTweetUrl(tweet) + "\">" + tweet.User.Name + "</a>"
	}

	// Go through each URL object and replace it with a link and correct text
	urls := []string{}
	for _, url := range tweet.Entities.Urls {
		replacement := "<a href='" + url.ExpandedURL + "'>" + url.DisplayURL + "</a>"
		from := url.Indices[0]
		to := url.Indices[1]
		// In case a tweet is shared with a comment there is a twitter URL in it
		// which we don't want to have since we will get the shared tweet later
		if tweet.QuotedStatus != nil {
			// This if does not one specific cases when a user RT (id1) a
			// shared post (id2) that shares the final post (id3). This
			// checks id1 == id2 and not id1 == id3
			if strings.EqualFold(url.ExpandedURL, GetTweetUrl(*tweet.QuotedStatus)) {
				// replace with nothing
				replacement = ""
			}
		}
		replacements = append(replacements, replaceObject{from, to, replacement})
		if len(replacement) > 0 {
			urls = append(urls, url.ExpandedURL)
		}
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
		text += "\n<blockquote>\n" + ParseTweetText(*tweet.QuotedStatus) + "\n</blockquote>"
	}

	footer := "\n<p><a href=\"" + GetTweetUrl(tweet) + "\">" + tweet.User.Name + "</a> @ " + ConvertTwitterTime(tweet.CreatedAt).Format(time.RFC1123) + "\n"

	// process all urls we found so we can append it to the entry
	var urlText string
	if len(urls) > 0 {
		urlText = "<hr>\n" + getUrlContent(urls) + "\n"
	}

	return text + footer + urlText
}
