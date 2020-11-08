package cyoa

import (
	"html/template"
	"net/http"
	"net/url"
)

type HandlerOption func(h *storyHandler)

// functional option. pretty cool.
func WithTemplate(tmpl *template.Template) HandlerOption {
	return func(h *storyHandler) {
		h.tmpl = tmpl
	}
}

func NewHandler(s Story, opts ...HandlerOption) http.Handler {
	h := storyHandler{story: s}
	for _, opt := range opts {
		opt(&h)
	}
	return h
}

type storyHandler struct {
	story map[string]StoryArc
	tmpl  *template.Template
}

func (sh storyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	u, _ := url.Parse(r.URL.Path)
	if u.Path == "/" {
		u.Path = "/intro"
		http.Redirect(w, r, u.String(), http.StatusFound)
	}
	arc, ok := sh.story[u.Path[1:]]
	if ok != true {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	} else {
		sh.tmpl.Execute(w, arc)
	}
}
