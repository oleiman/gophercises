package sitemapper

import (
	"encoding/xml"
	"fmt"
	"net/url"
	"os"
	// "gophercises/pkg/linkparser"
)

type URL struct {
	Location string `xml:"loc"`
}

type SiteMap struct {
	XMLName xml.Name `xml:"urlset"`
	Xmlns   string   `xml:"xmlns,attr"`
	Urls    []URL    `xml:"url"`
}

func getPath(domain string, href string) ([]byte, error) {
	tmp, err := url.Parse(href)
	if err != nil {
		return nil, err
	} else if len(tmp.Host) == 0 || tmp.Host == domain {
		return []byte(tmp.Path), nil
	}
	return nil, fmt.Errorf("Domain mismatch: %s", href)

}

func bfs(u *url.URL) []URL {
	urls := []URL{}
	queue := []string{u.Path}
	visited := map[string]bool{
		u.Path: true,
	}

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		p, err := getPath(u.Host, curr)
		if err != nil {
			// fmt.Fprintf(os.Stderr, "%s\n", err)
			continue
		}
		u.Path = string(p)
		links, err := fetchLinks(u.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			continue
		}

		urls = append(urls, URL{u.String()})
		for _, link := range links {
			if visited[link.Href] == false {
				queue = append(queue, link.Href)
				visited[link.Href] = true
			}
		}
	}

	return urls
}

func GenerateSiteMap(rawurl string) ([]byte, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	urls := bfs(u)

	sitemap := SiteMap{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		Urls:  urls,
	}

	output, err := xml.MarshalIndent(sitemap, "  ", "    ")
	if err != nil {
		return nil, err
	}
	return []byte(xml.Header + string(output)), nil
}
