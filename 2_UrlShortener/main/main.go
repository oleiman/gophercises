package main

import (
	"flag"
	"fmt"
	"github.com/boltdb/bolt"
	"io/ioutil"
	"net/http"
	"urlshort"
)

func main() {
	mux := defaultMux()

	var yml string
	flag.StringVar(&yml, "yml",
		"./students/baltuky/redirect.yaml",
		"YAML spec for URL mapping")
	var json string
	flag.StringVar(&json, "json",
		"./redirect.json",
		"JSON spec for URL mapping")
	var db string
	flag.StringVar(&db, "db",
		"./redirect.db",
		"Bolt DB instance containing URL mappings")

	flag.Parse()

	ymlData, err := ioutil.ReadFile(yml)
	if err != nil {
		fmt.Print("YAML ERR: ")
		fmt.Println(err)
		ymlData = []byte("")
	}

	jsonData, err := ioutil.ReadFile(json)
	if err != nil {
		fmt.Print("JSON ERR: ")
		fmt.Println(err)
		jsonData = []byte("")
	}

	blt, err := bolt.Open(db, 0600, nil)
	if err != nil {
		fmt.Print("BOLT ERR: ")
		fmt.Println(err)
		panic(err)
	}

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	fmt.Print("IN MEMORY: ")
	fmt.Println(pathsToUrls)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	// yaml := `
	// - path: /urlshort
	//   url: https://github.com/gophercises/urlshort
	// - path: /urlshort-final
	//   url: https://github.com/gophercises/urlshort/tree/solution
	// `
	yamlHandler, err := urlshort.YAMLHandler(ymlData, mapHandler)
	if err != nil {
		panic(err)
	}

	jsonHandler, err := urlshort.JSONHandler(jsonData, yamlHandler)
	if err != nil {
		panic(err)
	}

	boltHandler, err := urlshort.BoltHandler(blt, jsonHandler)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", boltHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
