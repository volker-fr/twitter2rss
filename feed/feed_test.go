package feed

import (
	"testing"

	"github.com/gorilla/feeds"
)

type headerTest struct {
	feed feeds.Feed
}

var headerTests = []headerTest{
	{
		feed: feeds.Feed{
			Title:       "Twitter Home Timeline",
			Description: "Twitter Home Timeline RSS Feed",
			Author:      &feeds.Author{Name: "Twitter2RSS", Email: "lists.volker@gmail.com"},
			Link:        &feeds.Link{Href: "http://github.com/volker-fr/twitter2rss/"},
			Items:       []*feeds.Item{},
		},
	},
}

func TestCreateFeedHeader(t *testing.T) {
	// Difficult to test since reflect.DeepEqual doesn't even work for an empty
	//  []*feeds.Item{}
	for _, testData := range headerTests {
		v := createFeedHeader()
		if testData.feed.Title != v.Title {
			t.Error("TestCreateFeedHeader title are different: %q != %q", testData.feed.Title, v.Title)
		}
		if testData.feed.Description != v.Description {
			t.Error("TestCreateFeedHeader description are different: %q != %q", testData.feed.Description, v.Description)
		}
	}
}
