package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"codeberg.org/textwire/textwire/v5"
	"codeberg.org/textwire/textwire/v5/config"
)

var Env = "development"

var tpl *textwire.Template
var host = "http://registry:5000"

type RepositoriesResponse struct {
	Repositories []string `json:"repositories"`
}

type ManifestResponse struct {
	Annotations struct {
		Arch string `json:"com.docker.official-images.bashbrew.arch,omitempty"`
	} `json:"annotations,omitempty"`
}

type Tag struct {
	Name string
	Arch string
}

type TagsResponse struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

func main() {
	isDev := os.Getenv("APP_ENV") == "development"

	twConf := &config.Config{
		DebugMode:     isDev,
		FileWatcher:   isDev,
		ErrorPagePath: "error-page",
		GlobalData: map[string]any{
			"isDev": isDev,
		},
	}

	t, err := textwire.NewTemplate(twConf)
	if err != nil {
		fmt.Println(err)
	}
	tpl = t

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	http.HandleFunc("/", reposHandler)
	http.HandleFunc("/tags/{name}", tagsHandler)
	http.ListenAndServe(":8080", nil)
}

func reposHandler(w http.ResponseWriter, r *http.Request) {
	url := host + "/v2/_catalog?n=1000"
	reposResp, err := fetch[RepositoriesResponse](url, http.MethodGet)
	if err != nil {
		log.Fatalln(err)
	}

	fail := tpl.Response(w, "views/repos", map[string]any{
		"title": "Repositories",
		"repos": reposResp.Repositories,
	})

	fail.PrintOnError()
}

func tagsHandler(w http.ResponseWriter, r *http.Request) {
	repoName := r.PathValue("name")

	url := fmt.Sprintf("%s/v2/%s/tags/list", host, repoName)
	resp, err := fetch[TagsResponse](url, http.MethodGet)
	if err != nil {
		log.Fatalln(err)
	}

	tags := make([]Tag, 0)

	for _, tagName := range resp.Tags {
		url := fmt.Sprintf("%s/v2/%s/manifests/%s", host, repoName, tagName)
		resp, err := fetch[ManifestResponse](url, http.MethodGet)
		if err != nil {
			log.Fatalln(err)
		}

		tags = append(tags, Tag{
			Name: tagName,
			Arch: resp.Annotations.Arch,
		})
	}

	fail := tpl.Response(w, "views/tags", map[string]any{
		"name": repoName,
		"tags": tags,
	})

	fail.PrintOnError()
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
	req.Header.Set("Accept", "application/vnd.oci.image.manifest.v1+json")

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
