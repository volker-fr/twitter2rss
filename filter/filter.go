package filter

import (
	"fmt"
	"regexp"

	"github.com/volker-fr/twitter2rss/config"

	"github.com/dghubble/go-twitter/twitter"
)

func IsTweetFiltered(tweet twitter.Tweet, conf config.Config, parsedTweetText string) bool {
	if len(conf.IgnoreSource) > 0 {
		for _, searchString := range conf.IgnoreSource {
			re := regexp.MustCompile(searchString)
			if re.MatchString(tweet.Source) {
				if conf.Debug {
					fmt.Printf("Source filter regex matched for tweet %s: %s\n", tweet.IDStr, searchString)
				}
				return true
			}
		}
	}

	if len(conf.IgnoreText) > 0 {
		for _, searchString := range conf.IgnoreText {
			re := regexp.MustCompile(searchString)
			if re.MatchString(parsedTweetText) {
				if conf.Debug {
					fmt.Printf("Text filter regex matched for tweet %s: %s\n", tweet.IDStr, searchString)
				}
				return true
			}
		}
	}

	return false
}
