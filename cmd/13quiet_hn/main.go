package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/gophercises/quiet_hn/hn"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"
)

type result struct {
	idx int
	item
	err error
}

type pair struct {
	idx int
	id  int
}

type Cache struct {
	entries map[int]*item
	mtx     *sync.Mutex
}

func hnWorker(client *hn.Client, cache *Cache, jobs <-chan pair, results chan<- result) {
	for pr := range jobs {
		// check for item in cache
		cache.mtx.Lock()
		if val, ok := cache.entries[pr.id]; ok {
			cache.mtx.Unlock()
			results <- result{idx: pr.idx, item: *val}
			continue
		}
		cache.mtx.Unlock()
		hnItem, err := client.GetItem(pr.id)
		if err != nil {
			results <- result{idx: pr.idx, err: err}
			continue
		}
		item := parseHNItem(hnItem)

		// add the item to the cache
		cache.mtx.Lock()
		cache.entries[pr.id] = &item
		cache.mtx.Unlock()

		if isStoryLink(item) {
			results <- result{idx: pr.idx, item: item}
		} else {
			results <- result{idx: pr.idx, err: errors.New("Not a story")}
		}
	}
}

func main() {
	// parse flags
	var port, numStories int
	flag.IntVar(&port, "port", 3000, "the port to start the web server on")
	flag.IntVar(&numStories, "num_stories", 30, "the number of top stories to display")
	flag.Parse()

	tpl := template.Must(template.ParseFiles("./index.gohtml"))

	cache := Cache{entries: make(map[int]*item), mtx: &sync.Mutex{}}
	http.HandleFunc("/", handler(numStories, tpl, &cache))

	// Start the servermake
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func handler(numStories int, tpl *template.Template, cache *Cache) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		var client hn.Client
		ids, err := client.TopItems()
		if err != nil {
			http.Error(w, "Failed to load top stories", http.StatusInternalServerError)
			return
		}
		var results []result
		nWorkers := 10
		want := numStories + 10
		jobsCh := make(chan pair, want)
		resultsCh := make(chan result, want)
		for w := 0; w < nWorkers; w++ {
			go hnWorker(&client, cache, jobsCh, resultsCh)
		}

		for i, id := range ids[:want] {
			jobsCh <- pair{idx: i, id: id}
		}
		close(jobsCh)
		for i := 0; i < want; i++ {
			res := <-resultsCh
			if res.err == nil {
				results = append(results, res)
			}
		}

		sort.Slice(results, func(i, j int) bool {
			return results[i].idx < results[j].idx
		})

		var stories []item
		for _, res := range results[:numStories] {
			stories = append(stories, res.item)
		}

		data := templateData{
			Stories: stories,
			Time:    time.Now().Sub(start),
		}
		err = tpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Failed to process the template", http.StatusInternalServerError)
			return
		}
	})
}

func isStoryLink(item item) bool {
	return item.Type == "story" && item.URL != ""
}

func parseHNItem(hnItem hn.Item) item {
	ret := item{Item: hnItem}
	url, err := url.Parse(ret.URL)
	if err == nil {
		ret.Host = strings.TrimPrefix(url.Hostname(), "www.")
	}
	return ret
}

// item is the same as the hn.Item, but adds the Host field
type item struct {
	hn.Item
	Host string
}

type templateData struct {
	Stories []item
	Time    time.Duration
}
