package urlshort

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"gopkg.in/yaml.v2"
	"net/http"
	"net/url"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		u, err := url.Parse(request.URL.Path)
		if err != nil {
			fallback.ServeHTTP(writer, request)
		}

		url := pathsToUrls[u.Path]
		if url == "" {
			fallback.ServeHTTP(writer, request)
		} else {
			// fmt.Fprintf(writer, "Hellosky %s!", request.URL.Path[1:])
			http.Redirect(writer, request, url, http.StatusFound)
		}
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.

func parseYAML(yml []byte) ([]map[string]string, error) {
	var mappingSeq []map[string]string
	err := yaml.Unmarshal(yml, &mappingSeq)
	return mappingSeq, err
}

func buildMap(ymlMaps []map[string]string) map[string]string {
	result := make(map[string]string)
	for _, mp := range ymlMaps {
		result[mp["path"]] = mp["url"]
	}
	return result
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	ymlMaps, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}
	pathsToUrls := buildMap(ymlMaps)
	fmt.Print("YAML: ")
	fmt.Println(pathsToUrls)
	return MapHandler(pathsToUrls, fallback), nil
}

func JSONHandler(jsn []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var pathsToUrls map[string]string
	err := json.Unmarshal(jsn, &pathsToUrls)
	if err != nil {
		return nil, err
	}
	fmt.Print("JSON: ")
	fmt.Println(pathsToUrls)
	return MapHandler(pathsToUrls, fallback), nil
}

func BoltHandler(db *bolt.DB, fallback http.Handler) (http.HandlerFunc, error) {
	return func(writer http.ResponseWriter, request *http.Request) {
		u, err := url.Parse(request.URL.Path)
		if err != nil {
			fallback.ServeHTTP(writer, request)
		}
		var url string
		err = db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("Mappings"))
			if b == nil {
				return errors.New("Bucket not found")
			}
			v := b.Get([]byte(u.Path))
			if v == nil {
				return errors.New("Path not found in DB")
			}
			url = string(v)
			return nil
		})

		if err != nil {
			fallback.ServeHTTP(writer, request)
		} else {
			http.Redirect(writer, request, url, http.StatusFound)
		}
	}, nil
}
