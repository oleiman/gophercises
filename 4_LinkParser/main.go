package main

import (
	"flag"
	"fmt"
	"linkParse/parser"
	"log"
)

func main() {
	var file string
	flag.StringVar(&file, "file", "ex1.html", "HTML file to parse")
	flag.Parse()

	doc, err := parser.ParseHtmlFromFile(file)
	if err != nil {
		log.Fatal(err)
	}
	links := parser.ExtractLinks(doc)
	fmt.Printf("%+v\n", links)
}
