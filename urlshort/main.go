package main

import (
	"fmt"
	"flag"
	"io/ioutil"
	"net/http"
)

import (
	"gopkg.in/yaml.v2"
)

func main() {
	yamlPath := flag.String("shorturls", "shorts.yaml", "Path to a YAML config for shortened URLs")
	flag.Parse()

	yaml, err := ioutil.ReadFile(*yamlPath)
	if err != nil {
		panic(err)
	}
	mux := defaultMux()
	pathMap := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := MapHandler(pathMap, mux)
	yamlHandler := YamlHandler([]byte(yaml), mapHandler)
	http.ListenAndServe(":8080", yamlHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, world!")
	})
	return mux
}

func MapHandler(pathMap map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if redirect, ok := pathMap[path]; ok {
			http.RedirectHandler(redirect, http.StatusFound).ServeHTTP(w, r)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

func YamlHandler(yml []byte, fallback http.Handler) http.HandlerFunc {
	type ShortURL struct {
		Path string
		URL string `yaml:"url"`
	}
	var urls []ShortURL
	err := yaml.Unmarshal(yml, &urls)
	if err != nil {
		return fallback.ServeHTTP
	}
	pathMap := make(map[string]string, len(urls))
	for _, shortURL := range urls {
		pathMap[shortURL.Path] = shortURL.URL
	}
	return MapHandler(pathMap, fallback)
}
