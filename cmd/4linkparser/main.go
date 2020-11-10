package main

import (
	"flag"
	"fmt"
	"gophercises/pkg/linkparser"
	"log"
)

func main() {
	var file string
	flag.StringVar(&file, "file", "data/4/ex1.html", "HTML file to parse")
	flag.Parse()

	doc, err := linkparser.ParseHtmlFromFile(file)
	if err != nil {
		log.Fatal(err)
	}
	links := linkparser.ExtractLinks(doc)
	fmt.Printf("%+v\n", links)
}
