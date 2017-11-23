package main

import (
	"encoding/json"
	"flag"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
)

type Scene struct {
	Title string
	Story []string
	Options []SceneOption
}

func (scene *Scene) StoryBlob() string {
	return strings.Join(scene.Story, " ")
}

type SceneOption struct {
	Text string
	Arc string
}

func buildHandler(scenes map[string]Scene, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		arc := r.URL.Path[1:]
		if arc == "" {
			http.Redirect(w, r, "/intro", http.StatusMovedPermanently)
			return
		}
		scene, ok := scenes[arc]
		if !ok {
			http.NotFound(w, r)
			return
		}
		err := tmpl.Execute(w, &scene)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func scenes(storyFilePath string) (map[string]Scene, error) {
	raw, err := ioutil.ReadFile(storyFilePath)
	if err != nil {
		return nil, err
	}
	sceneMap := make(map[string]Scene)
	err = json.Unmarshal([]byte(raw), &sceneMap)
	if err != nil {
		return nil, err
	}
	return sceneMap, nil
}

func main() {
	storyFilePath := flag.String("story-file", "gopher.json", "Path to story in JSON")
	flag.Parse()

	sceneMap, err := scenes(*storyFilePath)
	if err != nil {
		panic(err)
	}
	sceneTmpl := template.Must(template.ParseFiles("scene.html.tmpl"))
	handler := buildHandler(sceneMap, sceneTmpl)
	http.Handle("/", handler)
	http.ListenAndServe(":8080", handler)
}
