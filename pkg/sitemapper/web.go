package sitemapper

import (
	"gophercises/pkg/linkparser"
	"net/http"
)

// TODO(oren): reading the response body and returning
// only to create a new reader later and read them again is
// a bunch of extra work. not great.
func fetchLinks(url string) ([]linkparser.Link, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	doc, err := linkparser.ParseHtml(resp.Body)
	if err != nil {
		return nil, err
	}

	links := linkparser.ExtractLinks(doc)

	return links, nil

}
