package main

import (
	"adventure/cyoa"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	var fname string
	flag.StringVar(&fname, "story", "gopher.json", "JSON file containing a CYOA story")

	flag.Parse()

	storyData, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Fatal(err)
	}

	story, err := cyoa.ParseStory(storyData)
	if err != nil {
		log.Fatal(err)
	}

	tmpl, err := template.ParseFiles("layout.html")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Starting server on :8080")
	http.ListenAndServe(":8080", cyoa.NewHandler(story, cyoa.WithTemplate(tmpl)))
}
