package filter

import (
	"testing"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/volker-fr/twitter2rss/config"
)

type testSet struct {
	tweet      twitter.Tweet
	config     config.Config
	parsedText string
	filtered   bool
}

var tweets = []testSet{
	{
		tweet:      twitter.Tweet{Source: "first"},
		parsedText: "first",
		config: config.Config{IgnoreText: []string{"abc", "def"},
			IgnoreSource: []string{"ghi", "jkl"}},
		filtered: false,
	},
	{
		tweet:      twitter.Tweet{Source: "second"},
		parsedText: "second",
		config: config.Config{IgnoreText: []string{"abc", "def", "eco"},
			IgnoreSource: []string{"ghi", "jkl"}},
		filtered: true,
	},
	{
		tweet:      twitter.Tweet{Source: "third"},
		parsedText: "third",
		config: config.Config{IgnoreText: []string{"abc", "def"},
			IgnoreSource: []string{"ghi", "jkl", "ird"}},
		filtered: true,
	},
}

func TestIsTweetFiltered(t *testing.T) {
	for _, testData := range tweets {
		v := IsTweetFiltered(testData.tweet, testData.config, testData.parsedText)
		if v != testData.filtered {
			t.Error("Expected the filter to cache test tweet", testData.parsedText)
		}
	}
}
