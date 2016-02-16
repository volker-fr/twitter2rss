package feed

import (
	"reflect"
	"testing"

	"github.com/dghubble/go-twitter/twitter"
)

type testSet struct {
	tweets      []twitter.Tweet
	segmentSize int
	sorted      map[string]map[string][]twitter.Tweet
}

var tweets = []testSet{
	{
		segmentSize: 6,
		tweets: []twitter.Tweet{
			twitter.Tweet{User: &twitter.User{ScreenName: "a"}, CreatedAt: "Fri Nov 04 00:00:00 +0000 2011"},
			twitter.Tweet{User: &twitter.User{ScreenName: "a"}, CreatedAt: "Fri Nov 04 01:00:00 +0000 2011"},
			twitter.Tweet{User: &twitter.User{ScreenName: "a"}, CreatedAt: "Fri Nov 04 07:00:00 +0000 2011"},
			twitter.Tweet{User: &twitter.User{ScreenName: "b"}, CreatedAt: "Fri Nov 04 21:22:36 +0000 2011"},
		},
		sorted: map[string]map[string][]twitter.Tweet{
			"2011-11-04-0": {
				"a": {
					twitter.Tweet{User: &twitter.User{ScreenName: "a"}, CreatedAt: "Fri Nov 04 00:00:00 +0000 2011"},
					twitter.Tweet{User: &twitter.User{ScreenName: "a"}, CreatedAt: "Fri Nov 04 01:00:00 +0000 2011"},
				},
			},
			"2011-11-04-1": {
				"a": {
					twitter.Tweet{User: &twitter.User{ScreenName: "a"}, CreatedAt: "Fri Nov 04 07:00:00 +0000 2011"},
				},
			},
			"2011-11-04-3": {
				"b": {
					twitter.Tweet{User: &twitter.User{ScreenName: "b"}, CreatedAt: "Fri Nov 04 21:22:36 +0000 2011"},
				},
			},
		},
	},
}

func TestSortTweetsIntoHourSegments(t *testing.T) {
	for _, testData := range tweets {
		v := sortTweetsIntoHourSegments(testData.tweets, testData.segmentSize)
		// TODO: find better comparison as reflect.DeepEqual. The order of
		//       the input array matters with reflect.DeepEqual
		eq := reflect.DeepEqual(testData.sorted, v)
		if !eq {
			t.Error("Expected that output is correctly sorted by sortTweetsIntoHourSegments")
		}
	}
}
