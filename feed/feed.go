package feed

import (
	"fmt"
	"strconv"
	"time"

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
		Link:        &feeds.Link{Href: "http://github.com/volker-fr/twitter2rss/"},
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

		titleLimit := 60
		if len(tweet.Text) < 60 {
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

	sortedTweets := sortTweetsIntoHourSegments(tweets, conf.CombinedFeedHours)

	for tweetTimeSegment, timeSortedTweets := range sortedTweets {
		// get YYYY-MM-DD + Hour / <segmentsize>
		hourBlock := strconv.Itoa(time.Now().Hour() / conf.CombinedFeedHours)
		currentTimeSegment := time.Now().Format("2006-01-02") + "-" + hourBlock

		if tweetTimeSegment == currentTimeSegment {
			fmt.Println("INFO: Skipping tweets from the most recent time segment.")
			fmt.Println("      This will avoid duplicates or incomplete rss entries.")
			continue
		}

		for twitterUser, authorTweets := range timeSortedTweets {
			var feedText string

			for _, tweet := range timeSort(authorTweets) {
				parsedTweetText := parser.ParseTweetText(tweet)

				if filter.IsTweetFiltered(tweet, conf, parsedTweetText) {
					continue
				}
				feedText += parsedTweetText + "\n<hr><hr>\n"
			}

			// Calculate the time so we have a nicer formating
			segment, err := strconv.Atoi(tweetTimeSegment[11:])
			if err != nil {
				fmt.Printf("WARNING: couldn't converte %q to integer\n", tweetTimeSegment[11:])
			}
			segmentTime := strconv.Itoa(conf.CombinedFeedHours * segment)

			rssDate, err := time.Parse("2006-01-02-15", tweetTimeSegment[0:10]+"-"+segmentTime)
			if err != nil {
				fmt.Printf("WARNING: couldn't parse time %q \n", tweetTimeSegment[0:10]+"-"+segmentTime)
			}
			humanTweetSummaryTime := rssDate.Format("03pm")

			item := &feeds.Item{
				// TODO: check if slicing a string with non ascii chars will fail/scramble the text
				Title:       fmt.Sprintf("%s %s %s...", twitterUser, humanTweetSummaryTime, "combined tweets update"),
				Link:        &feeds.Link{Href: "https://twitter.com/" + twitterUser + "/"},
				Description: feedText,
				Author:      &feeds.Author{Name: twitterUser, Email: twitterUser},
				Created:     rssDate,
				Id:          fmt.Sprintf("%s-%s", currentTimeSegment, twitterUser),
			}
			feed.Add(item)
		}
	}

	return feed
}
