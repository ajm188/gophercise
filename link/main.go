package main

import (
	"fmt"
	"io"
	"os"
)

import (
	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

func findLinks(r io.Reader) []Link {
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
				links = append(links, Link{href, text})
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

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <file>\n", os.Args[0])
		os.Exit(1)
	}
	filepath := os.Args[1]
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	fmt.Println(findLinks(file))
}
