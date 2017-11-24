package main

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sync"
)

import (
	"github.com/ajm188/gophercise/link/link"
)

const (
	DEFAULT_SITEMAP_XMLNS = "http://www.sitemaps.org/schemas/sitemap/0.9"
)

type Sitemap struct {
	Urlsets []URLSet
}

func (sitemap *Sitemap) AddURLset() {
	urlset := URLSet{XMLNS: DEFAULT_SITEMAP_XMLNS}
	sitemap.Urlsets = append(sitemap.Urlsets, urlset)
}

func (sitemap *Sitemap) AddLink(lnk *url.URL) {
	sitemap.Urlsets[0].AddLink(lnk)
}

type URLSet struct {
	XMLName xml.Name `xml:"urlset"`
	XMLNS   string   `xml:"xmlns,attr"`
	Urls    []URL    `xml:"url"`
}

func (urlset *URLSet) AddLink(lnk *url.URL) {
	if urlset.Urls == nil {
		urlset.Urls = make([]URL, 0, 1)
	}
	urlset.Urls = append(urlset.Urls, URL{lnk.String()})
}

type URL struct {
	Loc string `xml:"loc"`
}

func crawl(start *link.Link) []link.Link {
	base := start.Href
	domain := start.Href.Hostname()
	fmt.Println(domain)
	cache := make(map[string]link.Link)
	mux := sync.Mutex{}

	_crawl := func(url link.Link, queue chan link.Link, done chan bool) {
		urlString := url.Href.String()
		defer func() { done <- true }()
		mux.Lock()
		_, visited := cache[urlString]
		mux.Unlock()
		if visited {
			return
		}
		fmt.Println("crawling ", url.Href)

		mux.Lock()
		cache[urlString] = url
		mux.Unlock()
		resp, err := http.Get(urlString)
		if err != nil {
			return
		}

		for _, lnk := range link.FindLinks(resp.Body) {
			href := lnk.Href
			lnk.Href = (*base).ResolveReference(href)
			if lnk.Href.Hostname() != domain {
				continue
			}
			queue <- lnk
		}
	}

	queue := make(chan link.Link)
	done := make(chan bool)
	defer close(queue)
	defer close(done)
	go _crawl(*start, queue, done)
	crawlers := 1
	for crawlers != 0 {
		select {
		case lnk := <-queue:
			go _crawl(lnk, queue, done)
			crawlers++
			fmt.Println("crawlers: ", crawlers)
		case _ = <-done:
			crawlers--
			fmt.Println("crawlers: ", crawlers)
		default:
			continue
		}
	}

	result := make([]link.Link, 0, len(cache))
	for _, v := range cache {
		result = append(result, v)
	}
	return result
}

func buildSitemap(links []link.Link) *Sitemap {
	sitemap := &Sitemap{}
	sitemap.AddURLset()
	for _, lnk := range links {
		sitemap.AddLink(lnk.Href)
	}
	return sitemap
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <url>\n", os.Args[0])
		os.Exit(1)
	}
	start, err := link.NewLink(os.Args[1], "")
	if err != nil {
		panic(err)
	}
	links := crawl(start)
	sitemap := buildSitemap(links)
	var out []byte
	for _, urlset := range sitemap.Urlsets {
		bytes, err := xml.Marshal(urlset)
		if err != nil {
			panic(err)
		}
		out = append(out, bytes...)
	}
	fmt.Printf("%s%s\n", xml.Header, out)
}
