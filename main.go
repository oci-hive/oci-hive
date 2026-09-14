package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"codeberg.org/textwire/textwire/v5"
	"codeberg.org/textwire/textwire/v5/config"
)

var tpl *textwire.Template
var host = "https://registry.serhiicho.com/"

func main() {
	tpl, _ = textwire.NewTemplate(&config.Config{
		DebugMode:   true,
		FileWatcher: true,
	})

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", reposHandler)
	http.HandleFunc("/tags/{name}", tagsHandler)
	http.ListenAndServe(":8080", nil)
}

func reposHandler(w http.ResponseWriter, r *http.Request) {
	url := host + "/v2/_catalog?n=1000"
	resp, err := fetch[RepositoriesResponse](url, http.MethodGet)
	if err != nil {
		log.Fatalln(err)
	}

	fail := tpl.Response(w, "views/repos", map[string]any{
		"title": "Repositories",
		"repos": resp.Repositories,
	})

	fail.PrintOnError()
}

func tagsHandler(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")

	url := fmt.Sprintf("%s/v2/%s/tags/list", host, name)
	resp, err := fetch[TagsResponse](url, http.MethodGet)
	if err != nil {
		log.Fatalln(err)
	}

	fail := tpl.Response(w, "views/tags", map[string]any{
		"name": name,
		"tags": resp.Tags,
	})

	fail.PrintOnError()
}

type RepositoriesResponse struct {
	Repositories []string `json:"repositories"`
}

type TagsResponse struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

func fetch[T any](url, method string) (*T, error) {
	// Create a context with timeout (prevents hanging requests)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Build the request
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("request error: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Use a custom client with timeout
	client := &http.Client{Timeout: 10 * time.Second}

	// Execute the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do error: %s", err)
	}

	// Handle non-2xx status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status %d: %s\n", resp.StatusCode, body)
	}

	var result T
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode error: %s", err)
	}

	defer resp.Body.Close()

	return &result, nil
}
