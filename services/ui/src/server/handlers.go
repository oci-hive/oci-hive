package server

import (
	"fmt"
	"log"
	"net/http"
)

// var host = "http://registry:5000"
var host = "https://registry.serhiicho.com"

func reposHandler(w http.ResponseWriter, r *http.Request) {
	url := host + "/v2/_catalog?n=1000"
	reposResp, err := fetch[repositoriesResponse](url, http.MethodGet)
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
	resp, err := fetch[tagsResponse](url, http.MethodGet)
	if err != nil {
		log.Fatalln(err)
	}

	tags := make([]tag, 0)

	for _, tagName := range resp.Tags {
		url := fmt.Sprintf("%s/v2/%s/manifests/%s", host, repoName, tagName)
		resp, err := fetch[manifestResponse](url, http.MethodGet)
		if err != nil {
			log.Fatalln(err)
		}

		tags = append(tags, tag{
			Name: tagName,
			Arch: getArchFromManifest(resp),
		})
	}

	fail := tpl.Response(w, "views/tags", map[string]any{
		"name": repoName,
		"tags": tags,
	})

	fail.PrintOnError()
}

func getArchFromManifest(resp *manifestResponse) []string {
	if resp.Annotations.Arch != "" {
		return []string{resp.Annotations.Arch}
	}

	if len(resp.Manifests) > 0 {
		res := make([]string, 0, 1)
		for _, manifest := range resp.Manifests {
			res = append(res, manifest.Platform.Arch)
		}
		return res
	}

	return []string{}
}
