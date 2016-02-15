package feed

import (
	"fmt"
	"time"
	"strconv"

	"github.com/volker-fr/twitter2rss/config"
	"github.com/volker-fr/twitter2rss/filter"
	"github.com/volker-fr/twitter2rss/parser"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/gorilla/feeds"
)

func createFeedHeader() *feeds.Feed {
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

	return feed
}

// Every tweet is its own feed item
func CreateIndividualFeed(conf config.Config, tweets []twitter.Tweet) *feeds.Feed {
	feed := createFeedHeader()

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

	return feed
}

// combine multiple tweets together into one feed item
func CreateCombinedUserFeed(conf config.Config, tweets []twitter.Tweet) *feeds.Feed {
	feed := createFeedHeader()

	// TODO: move
	segmentSize := 6

	sortedTweets := sortTweetsIntoHourSegments(tweets, segmentSize)

	for timeSegment, timeSortedTweets := range sortedTweets {
		// get YYYY-MM-DD + Hour / <segmentsize>
		hourBlock := strconv.Itoa(time.Now().Hour() / segmentSize)
		currentSegment := time.Now().Format("2006-01-02") + "-" + hourBlock

		if timeSegment == currentSegment {
			fmt.Printf("Skipping tweets from the last segment...")
			continue
		}

		for twitterUser, authorTweets := range timeSortedTweets {
			var feedText string

			for _, tweet := range authorTweets {
				parsedTweetText := parser.ParseTweetText(tweet)

				if filter.IsTweetFiltered(tweet, conf, parsedTweetText) {
					continue
				}
				feedText += parsedTweetText + "\n<hr><hr>\n"
			}

			// the last value is unrelated and not a real time but the internal
			// time segment, only 2006-01-02 matters
			rssDate, _ := time.Parse("2006-01-02-03", currentSegment)

			item := &feeds.Item{
				// TODO: check if slicing a string with non ascii chars will fail/scramble the text
				Title:       fmt.Sprintf("%s %s: %s...", currentSegment, twitterUser, "combined tweets"),
				Link:        &feeds.Link{Href: "https://twitter.com/" + twitterUser + "/"},
				Description: feedText,
				Author:      &feeds.Author{Name: twitterUser, Email: twitterUser},
				Created:     rssDate,
				Id:          fmt.Sprintf("%s-%s", currentSegment, twitterUser),
			}
			feed.Add(item)
		}
	}

	return feed
}
