package parser

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

type Content struct {
	URL				   string
	HTMLTitle          string
	TwitterTitle       string
	OGTitle            string
	HDLTitle           string
	Description        string
	OGDescription      string
	LPDescription      string
	TwitterDescription string
}

func getUrlContent(urls []string) string {
	// TODO: figure out what to do with non-HTML content such as PDF, images, ...
	var returnText string

	for _, url := range urls {
		fmt.Println(url)
		var content Content
		content.URL = url

		// open connection & get data
		resp, err := http.Get(url)
		defer resp.Body.Close()
		if err != nil {
			fmt.Printf("Couldn't load %s: %s", url, err)
			continue
		}

		z := html.NewTokenizer(resp.Body)
		for {
			tt := z.Next()
			// End of the document
			if tt == html.ErrorToken {
				break
			}
			if tt == html.StartTagToken {
				tagName, hasAttr := z.TagName()

				// parse title tag
				if strings.ToLower(string(tagName)) == "title" {
					// if its a title tag the next entry will be the text
					// TODO: check if this works with <title>a<b>c</b>d</title>
					z.Next()
					content.HTMLTitle = strings.Trim(z.Token().String(), " \t\n\r")
				} // if title

				// Parse meta tags
				if strings.ToLower(string(tagName)) == "meta" {
					var bKey, bVal []byte
					keys := make(map[string]string)
					// repeat until z.TagAttr reports that there are no more
					// attributes left (hasAttr = false)
					for hasAttr {
						bKey, bVal, hasAttr = z.TagAttr()
						keys[string(bKey)] = string(bVal)
					}
					content = getMetaInformation(content, keys)
				} // if meta
			} // if html.StartTagToken

		} // for html.NewTokenizer

		returnText += buildHTMLblock(content)
	} // for _, url
	return returnText
}


func getMetaInformation(content Content, keys map[string]string) Content {
	// Assign now all the various values we will find in keys to the content object

	// Description: <meta itemprop="description" name="description" content="...">
	if keys["name"] == "description" {
		content.Description = keys["content"]
	}
	// LPDescription: <meta name="lp" content="..." />
	if keys["name"] == "lp" {
		content.LPDescription = keys["content"]
	}
	// TwitterDescription: <meta property="twitter:description" content="..." />
	if keys["property"] == "twitter:description" || keys["name"] == "twitter:description" {
		content.TwitterDescription = keys["content"]
	}
	// OGDescription: <meta content="..." property="og:description" />
	if keys["property"] == "og:description" {
		content.OGDescription = keys["content"]
	}
	// TwitterTitle: <meta content="..." property="twitter:title" />
	if keys["property"] == "twitter:title" || keys["name"] == "twitter:title" {
		content.TwitterTitle = keys["content"]
	}
	// OGTitle: <meta content="..." property="og:title" />
	if keys["property"] == "og:title" {
		content.OGTitle = keys["content"]
	}
	// HDLTitle: <meta name="hdl" content"..." />
	if keys["name"] == "hdl" {
		content.HDLTitle = keys["content"]
	}

	return content
}


func buildHTMLblock(content Content) string {
	var title, description string

	// get the title
	if len(content.TwitterTitle) > 0 {
		title = content.TwitterTitle
	} else if len(content.OGTitle) > 0 {
		title = content.OGTitle
	} else if len(content.HDLTitle) > 0 {
		title = content.HDLTitle
	} else {
		title = content.HTMLTitle
	}

	// get the description
	if len(content.TwitterDescription) > 0 {
		description = content.TwitterDescription
	} else if len(content.OGDescription) > 0 {
		description = content.OGDescription
	} else if len(content.LPDescription) > 0 {
		description = content.LPDescription
	} else {
		description = content.Description
	}

	title = "<b>" + title + "</b><br>\n"
	footer := "<p><a href=\"" + content.URL + "\">" + content.URL + "</a>\n"

	return "<blockquote>\n" + title + description + footer + "</blockquote>\n"
}
