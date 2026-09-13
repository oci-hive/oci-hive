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

	http.HandleFunc("/", reposHandler)
	http.HandleFunc("/tags/{name}", tagsHandler)
	http.ListenAndServe(":8080", nil)
}

func reposHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := fetchRepositories()
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

	resp, err := fetchTags(name)
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

func fetchRepositories() (*RepositoriesResponse, error) {
	url := host + "/v2/_catalog?n=1000"

	// Create a context with timeout (prevents hanging requests)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Build the request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
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

	var result RepositoriesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode error: %s", err)
	}

	defer resp.Body.Close()

	return &result, nil
}

func fetchTags(name string) (*TagsResponse, error) {
	url := host + fmt.Sprintf("/v2/%s/tags/list", name)

	// Create a context with timeout (prevents hanging requests)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Build the request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
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

	var result TagsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode error: %s", err)
	}

	defer resp.Body.Close()

	return &result, nil
}
