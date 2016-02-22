package feed

import (
	"sort"
	"strconv"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/volker-fr/twitter2rss/parser"
)

type Tweets []twitter.Tweet

func (t Tweets) Len() int {
	return len(t)
}
func (t Tweets) Less(i, j int) bool {
	t1, _ := time.Parse(time.RubyDate, t[i].CreatedAt)
	t2, _ := time.Parse(time.RubyDate, t[j].CreatedAt)
	return t1.Before(t2)
}
func (t Tweets) Swap(i, j int) {
	//t[i], t[j] = t[j], t[i]
	temp := t[i]
	t[i] = t[j]
	t[j] = temp
}

func timeSort(tweets Tweets) Tweets {
	sort.Sort(tweets)
	return tweets
}

// sort tweet into hour segments & author
// Example:
// 		segment of  6 = 4x  6hour blocks a day
//	 	segment of 12 = 2x 12hour blocks a day
func sortTweetsIntoHourSegments(tweets []twitter.Tweet, segmentSize int) map[string]map[string]Tweets {
	// first string is the date, second the author
	var sortedList map[string]map[string]Tweets
	sortedList = make(map[string]map[string]Tweets)
	for _, tweet := range tweets {
		author := tweet.User.ScreenName
		tweetDate := parser.ConvertTwitterTime(tweet.CreatedAt)
		hourBlock := strconv.Itoa(tweetDate.Hour() / segmentSize)
		timeSegment := tweetDate.Format("2006-01-02") + "-" + hourBlock

		if sortedList[timeSegment] == nil {
			sortedList[timeSegment] = make(map[string]Tweets)
		}
		sortedList[timeSegment][author] = append(sortedList[timeSegment][author], tweet)
	}

	return sortedList
}
