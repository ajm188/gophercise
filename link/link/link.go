package link

import (
	"io"
	"net/url"
)

import (
	"golang.org/x/net/html"
)

type Link struct {
	Href *url.URL
	Text string
}

func NewLink(href, text string) (*Link, error) {
	parsedHref, err := url.Parse(href)
	if err != nil {
		return nil, err
	}
	return &Link{parsedHref, text}, nil
}

func FindLinks(r io.Reader) []Link {
	tokenizer := html.NewTokenizer(r)
	links := make([]Link, 0)
	for {
		switch tokenizer.Next() {
		case html.ErrorToken:
			return links
		case html.StartTagToken:
			tag, attrs := tokenizer.TagName()
			if string(tag) == "a" && attrs {
				href := findHref(tokenizer)
				text := findText(tokenizer)
				link, err := NewLink(href, text)
				if err != nil {
					panic(err)
				}
				links = append(links, *link)
			}
		}
	}
	return links
}

func findText(tokenizer *html.Tokenizer) string {
	fullText := ""
	for {
		switch tokenizer.Next() {
		case html.TextToken:
			fullText += string(tokenizer.Text())
		case html.EndTagToken:
			tag, _ := tokenizer.TagName()
			if string(tag) == "a" {
				return fullText
			}
		}
	}
	return fullText
}

func findHref(tokenizer *html.Tokenizer) string {
	for {
		name, val, hasMore := tokenizer.TagAttr()
		if string(name) == "href" {
			return string(val)
		}
		if !hasMore {
			return ""
		}
	}
}
