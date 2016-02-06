package parser

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

// TODO: there is also "<link itemprop=embedURL ..." but its much more
//		 difficult to parse since the hight & width value can occure multiple
//		 times if you look at the youtube html code
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
	TwitterPlayer	string
	TwitterPlayerWidth string
	TwitterPlayerHeight string
	OGVideoURL string
	OGVideoWidth string
	OGVideoHeight string
}

func getUrlContent(urls []string) string {
	// TODO: figure out what to do with non-HTML content such as PDF, images, ...
	var returnText string

	for _, url := range urls {
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
			if tt == html.StartTagToken || tt == html.SelfClosingTagToken || tt == html.EndTagToken {
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
	// TwitterPlayer: <meta name="twitter:player" content="...">
	if keys["name"] == "twitter:player" {
		content.TwitterPlayer = keys["content"]
	}
	// TwitterPlayerWidth: <meta name="twitter:player:width" content="...">
	if keys["name"] == "twitter:player:width" {
		content.TwitterPlayerWidth = keys["content"]
	}
	// TwitterPlayerHeight: <meta name="twitter:player:height" content="...">
	if keys["name"] == "twitter:player:height" {
		content.TwitterPlayerHeight = keys["content"]
	}
	// OGVideoURL: <meta property="og:video:url" content="...">
	if keys["property"] == "og:video:url" {
		content.OGVideoURL = keys["content"]
	}
	// OGVideoWidth: <meta property="og:video:width" content="...">
	if keys["property"] == "og:video:width" {
		content.OGVideoWidth = keys["content"]
	}
	// OGVideoHeight: <meta property="og:video:height" content="...">
	if keys["property"] == "og:video:height" {
		content.OGVideoHeight = keys["content"]
	}

	return content
}


func buildHTMLblock(content Content) string {
	var title, description, mediaText string
	var media, mediaWidth, mediaHeight string

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

	// get the media/video
	if len(content.TwitterPlayer) > 0 {
		media = content.TwitterPlayer
	} else if len(content.OGVideoURL) > 0 {
		media = content.OGVideoURL
	}
	// get media width
	if len(content.TwitterPlayerWidth) > 0 {
		mediaWidth = content.TwitterPlayerWidth
	} else if len(content.OGVideoWidth) > 0 {
		mediaWidth = content.OGVideoWidth
	}

	// get media height
	if len(content.TwitterPlayerHeight) > 0 {
		mediaHeight = content.TwitterPlayerHeight
	} else if len(content.OGVideoHeight) > 0 {
		mediaHeight = content.OGVideoHeight
	}

	if len(media) > 0 && len(mediaWidth) >0 && len(mediaHeight) > 0 {
		mediaText = "\n<iframe width=" + mediaWidth + " height=" + mediaHeight + " src=" + media + " frameborder=0 allowfullscreen></iframe>"
	}
	title = "<b>" + title + "</b><br>\n"
	footer := "<p><a href=\"" + content.URL + "\">" + content.URL + "</a>\n"

	return "<blockquote>\n" + title + description + footer + mediaText + "</blockquote>\n"
}
