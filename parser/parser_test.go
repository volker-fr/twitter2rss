package parser

import (
	"github.com/dghubble/go-twitter/twitter"
	"strings"
	"testing"
	"time"
)

type testSet struct {
	tweet           twitter.Tweet
	resultUrl       string
	resultFeedTexts []string
	resultRFC822    string
}

var tweets = []testSet{
	{
		tweet: twitter.Tweet{User: &twitter.User{ScreenName: "username"},
			IDStr:     "12345",
			Text:      "tweet text",
			CreatedAt: "Wed Aug 27 13:08:45 +0000 2008"},
		resultUrl:       "https://twitter.com/username/status/12345",
		resultFeedTexts: []string{"tweet text", "username", "Wed, 27 Aug 2008 13:08"},
		resultRFC822:    "27 Aug 08 13:08 +0000",
	},
}

func TestGetTweetUrl(t *testing.T) {
	for _, testData := range tweets {
		v := GetTweetUrl(testData.tweet)
		if v != testData.resultUrl {
			t.Error(
				"Expected", testData.resultUrl,
				"got", v,
			)
		}
	}
}

func TestConvertTwitterTime(t *testing.T) {
	for _, testData := range tweets {
		v := ConvertTwitterTime(testData.tweet.CreatedAt)
		if v.Format(time.RFC822) != testData.resultRFC822 {
			t.Error(
				"Expected", testData.resultRFC822,
				"got", v.Format(time.RFC822),
			)
		}
	}
}

func TestParseTweetText(t *testing.T) {
	for _, testData := range tweets {
		v := ParseTweetText(testData.tweet)
		for _, resultString := range testData.resultFeedTexts {
			if !strings.Contains(v, resultString) {
				t.Error(
					"Expected to find", resultString,
					"in", v,
				)
			}
		}
	}
}
