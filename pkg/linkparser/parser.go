package linkparser

import (
	"bufio"
	"golang.org/x/net/html"
	"io"
	"os"
	"strings"
)

type Link struct {
	Href string
	Text string
}

func findHref(node *html.Node) string {
	for _, attr := range node.Attr {
		if attr.Key == "href" {
			return attr.Val
		}
	}
	return ""
}

func findText(node *html.Node) string {
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			return strings.TrimSpace(c.Data)
		}
	}
	return ""
}

func ExtractLinks(root *html.Node) []Link {
	links := []Link{}
	if root.Type == html.ElementNode && root.Data == "a" {
		link := Link{findHref(root), findText(root)}
		links = append(links, link)
	} else {
		for c := root.FirstChild; c != nil; c = c.NextSibling {
			links = append(links, ExtractLinks(c)...)
		}
	}
	return links
}

func ParseHtml(reader io.Reader) (*html.Node, error) {
	doc, err := html.Parse(reader)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

// TODO(oren): would rather this took a reader or something
// maybe even a byte slice, but it's time to move on
func ParseHtmlFromFile(fname string) (*html.Node, error) {
	fd, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	reader := bufio.NewReader(fd)
	doc, err := html.Parse(reader)
	if err != nil {
		return nil, err
	}
	return doc, nil
}
