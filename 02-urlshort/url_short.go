package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jmaeso/gophercises/02-urlshort/api"
	yaml "gopkg.in/yaml.v2"
)

func main() {
	yamlFilePath := flag.String("yaml", "shortened_urls.yaml", "yaml file where to load shortened urls")
	jsonFilePath := flag.String("json", "shortened_urls.json", "json file where to load shortened urls")
	flag.Parse()

	mux := api.DefaultMux()

	// Build the MapHandler using the mux as the fallback
	shortenedURLsMap := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := api.MapHandler(shortenedURLsMap, mux)

	yamlStore, err := parseYAML(*yamlFilePath)
	if err != nil {
		panic(err)
	}
	yamlHandler := api.YAMLHandler(yamlStore, mapHandler)

	jsonStore, err := parseJSON(*jsonFilePath)
	if err != nil {
		panic(err)
	}
	jsonHandler := api.JSONHandler(jsonStore, yamlHandler)

	fmt.Println("Starting the server on :8080")
	log.Fatal(http.ListenAndServe(":8080", jsonHandler))
}

func parseYAML(yamlFilePath string) ([]api.ShortenedURL, error) {
	yamlFile, err := os.Open(yamlFilePath)
	if err != nil {
		return nil, err
	}

	yamlDecoder := yaml.NewDecoder(yamlFile)

	shortenedURLs := make([]api.ShortenedURL, 0)

	if err := yamlDecoder.Decode(&shortenedURLs); err != nil {
		return nil, err
	}

	return shortenedURLs, nil
}

func parseJSON(jsonFilePath string) ([]api.ShortenedURL, error) {
	jsonFile, err := os.Open(jsonFilePath)
	if err != nil {
		return nil, err
	}

	jsonDecoder := json.NewDecoder(jsonFile)

	jsonPayload := api.ShortenedURLJSONPayload{}

	if err := jsonDecoder.Decode(&jsonPayload); err != nil {
		return nil, err
	}

	return jsonPayload.ShortenedURL, nil
}
