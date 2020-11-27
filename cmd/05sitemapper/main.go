package main

import (
	"flag"
	"gophercises/pkg/sitemapper"
	"log"
	"os"
)

func main() {

	var url string
	var outfile string
	flag.StringVar(&url, "url", "http://calhoun.io", "URL of site to map")
	flag.StringVar(&outfile, "out", "", "Output file")
	flag.Parse()

	fd, err := func(outfile string) (*os.File, error) {
		switch {
		case outfile == "":
			return os.Stdout, nil
		default:
			fd, err := os.Create(outfile)
			if err != nil {
				return nil, err
			}
			return fd, nil
		}
	}(outfile)

	if err != nil {
		log.Fatal(err)
	}

	sitemap, err := sitemapper.GenerateSiteMap(url)
	if err != nil {
		log.Fatal(err)
	}
	fd.Write(sitemap)
}
