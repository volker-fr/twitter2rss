package feed

import (
	"strconv"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/volker-fr/twitter2rss/parser"
)

// sort tweet into hour segments & author
// Example:
// 		segment of  6 = 4x  6hour blocks a day
//	 	segment of 12 = 2x 12hour blocks a day
func sortTweetsIntoHourSegments(tweets []twitter.Tweet, segmentSize int) map[string]map[string][]twitter.Tweet {
	// first string is the date, second the author
	var sortedList map[string]map[string][]twitter.Tweet
	sortedList = make(map[string]map[string][]twitter.Tweet)
	for _, tweet := range tweets {
		author := tweet.User.ScreenName
		tweetDate := parser.ConvertTwitterTime(tweet.CreatedAt)
		hourBlock := strconv.Itoa(tweetDate.Hour() / segmentSize)
		timeSegment := tweetDate.Format("2006-01-02") + "-" + hourBlock

		if sortedList[timeSegment] == nil {
			sortedList[timeSegment] = make(map[string][]twitter.Tweet)
		}
		sortedList[timeSegment][author] = append(sortedList[timeSegment][author], tweet)
	}

	return sortedList
}
