package cyoa

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.New("").Parse(defaultHandlerTemplat))
}

var defaultHandlerTemplat = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Choose your own Adventure</title>
</head>
<body>
<h1>{{.Title}}</h1>
{{range .Paragraphs}}
    <p>{{.}}</p>
{{end}}
<ul>
    {{range .Options}}
        <li><a href="/{{.Chapter}}">{{.Text}}</a></li>
    {{end}}
</ul>
</body>
</html>`

type HandlerOption func(h *handler)

func WithTemplate(t *template.Template) HandlerOption {
	return func(h *handler) {
		h.t = t
	}
}

func WithPathFun(fn func(r *http.Request) string) HandlerOption {
	return func(h *handler) {
		h.pathFun = fn
	}
}

func defaultPathFun(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	if path == "" || path == "/" {
		path = "/intro"
	}
	path = path[1:]
	return path
}

func NewHandler(s Story, opts ...HandlerOption) http.Handler {
	h := handler{s, tpl, defaultPathFun}
	for _, opt := range opts {
		opt(&h)
	}
	return h
}

type handler struct {
	s       Story
	t       *template.Template
	pathFun func(r *http.Request) string
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := h.pathFun(r)
	if chapter, ok := h.s[path]; ok {
		err := h.t.Execute(w, chapter)
		//x := tpl.Tree
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "Something went Wrong ...", http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "Chapter not found.", http.StatusNotFound)
}

// code to add two numbers

func JsonStoryParser(r io.Reader) (Story, error) {
	decodedFile := json.NewDecoder(r)
	var story Story
	if err := decodedFile.Decode(&story); err != nil {
		return nil, err
	}
	return story, nil
}

type Story map[string]Chapter

type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

type Option struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}
