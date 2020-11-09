package parser_test

import (
	"linkParse/parser"
	"testing"
)

var sample = `
<html>
<body>
  <h1>Hello!</h1>
  <a href="/other-page">A link to another page</a>
</body>
</html>
`

var ex2 = "../ex2.html"

var ex2Expected = []parser.Link{
	parser.Link{
		Href: "https://www.twitter.com/joncalhoun",
		Text: "Check me out on twitter",
	},
	parser.Link{
		Href: "https://github.com/gophercises",
		Text: "Gophercises is on",
	},
}

var ex4 = "../ex4.html"
var ex4Text = "dog cat"

func TestParseHtml(t *testing.T) {
	doc, err := parser.ParseHtmlFromFile(ex2)
	if err != nil {
		t.Errorf("Failed to parse %s: %s", ex2, err)
	} else if doc == nil {
		t.Errorf("Doc root unexpectedly nil")
	}
}

func TestExtractLinks(t *testing.T) {
	doc, _ := parser.ParseHtmlFromFile(ex2)
	links := parser.ExtractLinks(doc)
	if len(links) != 2 {
		t.Errorf("Expected %d links, got %d", len(ex2Expected), len(links))
	}
	for i, link := range links {
		if link != ex2Expected[i] {
			t.Errorf("Link %d: expected %v, got %v", i, ex2Expected[i], link)
		}
	}
}

func TestNoComments(t *testing.T) {
	doc, _ := parser.ParseHtmlFromFile(ex4)
	links := parser.ExtractLinks(doc)
	if links[0].Text != ex4Text {
		t.Errorf("Expected \"%s\", got \"%s\"", ex4Text, links[0].Text)
	}
}
